const { makeRequest } = require('../utils/httpClient');
const { validateTask, validateTaskUpdate } = require('../utils/validation');
const logger = require('../utils/logger');
const cache = require('../utils/cache');

/**
 * Get all tasks
 */
const getAllTasks = async (req, res) => {
  try {
    const status = req.query.status || '';
    const userId = req.query.userId || '';
    const cacheKey = status || userId ? `tasks:status:${status}:userId:${userId}` : 'tasks';

    let url = '/api/tasks';
    if (status || userId) {
      const params = new URLSearchParams();
      if (status) params.append('status', status);
      if (userId) params.append('userId', userId);
      url += '?' + params.toString();
    }

    const response = await makeRequest(url);

    // Return shape frontend expects { tasks, count } (unwrap if Go sent envelope)
    const tasks = response?.tasks ?? response?.data?.tasks ?? [];
    const count = response?.count ?? response?.data?.count ?? tasks.length;
    const result = { tasks, count };

    // Refresh cache when API has more data than cache (keeps tasks in sync with stats)
    const cached = cache.get(cacheKey);
    let cachedCount = 0;
    if (cached) {
      try {
        const p = JSON.parse(cached);
        cachedCount = p?.count ?? p?.tasks?.length ?? 0;
      } catch {
        /* ignore */
      }
    }
    if (!cached || count > cachedCount) {
      cache.set(cacheKey, JSON.stringify(result));
    }

    return res.json(result);
  } catch (error) {
    const statusCode = error.statusCode || (error.message.includes('ECONNREFUSED') ? 503 : 500);
    logger.logError(req.method, req.originalUrl, statusCode, 0, 'GetTasksError', error.message);
    res.status(statusCode).json({
      error: error.message || 'Failed to fetch tasks',
      code: error.code || 'BACKEND_ERROR'
    });
  }
};

/**
 * Get task by ID
 */
const getTaskById = async (req, res) => {
  try {
    const cacheKey = `tasks:${req.params.id}`;
    const cached = cache.get(cacheKey);
    if (cached) {
      return res.json(JSON.parse(cached));
    }

    const raw = await makeRequest(`/api/tasks/${req.params.id}`);
    const task = raw?.data ?? raw;
    cache.set(cacheKey, JSON.stringify(task));
    res.json(task);
  } catch (error) {
    const statusCode = error.statusCode || (error.message.includes('not found') ? 404 : 500);
    logger.logError(req.method, req.originalUrl, statusCode, 0, 'GetTaskByIdError', error.message);
    res.status(statusCode).json({
      error: error.message || 'Task not found',
      code: error.code || 'NOT_FOUND'
    });
  }
};

/**
 * Create a new task
 */
const createTask = async (req, res) => {
  try {
    // Validate input
    const validation = validateTask(req.body);
    if (!validation.isValid) {
      logger.logError(req.method, req.originalUrl, 400, 0, 'ValidationError', validation.errors.join(', '));
      return res.status(400).json({ 
        error: validation.errors.join(', ') 
      });
    }
    
    const raw = await makeRequest('/api/tasks', {
      method: 'POST',
      body: req.body
    });
    const task = raw?.data ?? raw;
    cache.invalidate('tasks');
    res.status(201).json(task);
  } catch (error) {
    const statusCode = error.statusCode || (error.message.includes('400') ? 400 : 500);
    const errorName = statusCode === 400 ? 'ValidationError' : 'CreateTaskError';
    logger.logError(req.method, req.originalUrl, statusCode, 0, errorName, error.message);
    res.status(statusCode).json({ error: error.message });
  }
};

/**
 * Update a task
 */
const updateTask = async (req, res) => {
  try {
    // Validate input
    const validation = validateTaskUpdate(req.body);
    if (!validation.isValid) {
      logger.logError(req.method, req.originalUrl, 400, 0, 'ValidationError', validation.errors.join(', '));
      return res.status(400).json({ 
        error: validation.errors.join(', ') 
      });
    }
    
    const raw = await makeRequest(`/api/tasks/${req.params.id}`, {
      method: 'PUT',
      body: req.body
    });
    const task = raw?.data ?? raw;
    cache.invalidate('tasks');
    res.json(task);
  } catch (error) {
    if (error.message.includes('404') || error.message.includes('not found')) {
      logger.logError(req.method, req.originalUrl, 404, 0, 'TaskNotFoundError', error.message);
      res.status(404).json({ error: 'Task not found' });
    } else {
      const statusCode = error.statusCode || (error.message.includes('400') ? 400 : 500);
      const errorName = statusCode === 400 ? 'ValidationError' : 'UpdateTaskError';
      logger.logError(req.method, req.originalUrl, statusCode, 0, errorName, error.message);
      res.status(statusCode).json({ error: error.message });
    }
  }
};

module.exports = {
  getAllTasks,
  getTaskById,
  createTask,
  updateTask
};
