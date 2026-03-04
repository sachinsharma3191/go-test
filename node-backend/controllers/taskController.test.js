const {
  getAllTasks,
  getTaskById,
  createTask,
  updateTask,
  deleteTask
} = require('./taskController');
const { makeRequest } = require('../utils/httpClient');
const { validateTask, validateTaskUpdate } = require('../utils/validation');
const logger = require('../utils/logger');
const cache = require('../utils/cache');

// Mock dependencies
jest.mock('../utils/httpClient');
jest.mock('../utils/validation');
jest.mock('../utils/logger');
jest.mock('../utils/cache');

describe('Task Controller', () => {
  let mockReq, mockRes;

  beforeEach(() => {
    jest.clearAllMocks();

    mockReq = {
      method: 'GET',
      originalUrl: '/api/tasks',
      params: {},
      body: {},
      query: {}
    };

    mockRes = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis()
    };
  });

  describe('getAllTasks', () => {
    it('should fetch tasks and update cache when backend has data', async () => {
      const backendResponse = {
        tasks: [
          { id: 1, title: 'Task 1', status: 'pending', userId: 1 },
          { id: 2, title: 'Task 2', status: 'completed', userId: 2 }
        ],
        count: 2
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('tasks');
      expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      expect(cache.set).toHaveBeenCalledWith('tasks', JSON.stringify(backendResponse));
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should fetch tasks when not cached', async () => {
      const backendResponse = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('tasks');
      expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      expect(cache.set).toHaveBeenCalledWith('tasks', JSON.stringify(backendResponse));
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should handle status filter query parameter', async () => {
      mockReq.query.status = 'pending';

      const backendResponse = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks?status=pending');
      expect(cache.set).toHaveBeenCalledWith('tasks:status:pending:userId:', JSON.stringify(backendResponse));
    });

    it('should handle userId filter query parameter', async () => {
      mockReq.query.userId = '1';

      const backendResponse = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks?userId=1');
      expect(cache.set).toHaveBeenCalledWith('tasks:status::userId:1', JSON.stringify(backendResponse));
    });

    it('should handle both status and userId filter parameters', async () => {
      mockReq.query.status = 'pending';
      mockReq.query.userId = '1';

      const backendResponse = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks?status=pending&userId=1');
      expect(cache.set).toHaveBeenCalledWith('tasks:status:pending:userId:1', JSON.stringify(backendResponse));
    });

    it('should refresh cache when backend has more data than cache', async () => {
      const cachedTasks = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      const backendResponse = {
        tasks: [
          { id: 1, title: 'Task 1', status: 'pending', userId: 1 },
          { id: 2, title: 'Task 2', status: 'completed', userId: 2 }
        ],
        count: 2
      };

      cache.get.mockReturnValue(JSON.stringify(cachedTasks));
      makeRequest.mockResolvedValue(backendResponse);

      await getAllTasks(mockReq, mockRes);

      expect(cache.set).toHaveBeenCalledWith('tasks', JSON.stringify(backendResponse));
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should keep cache when backend has same or less data', async () => {
      const cachedTasks = {
        tasks: [
          { id: 1, title: 'Task 1', status: 'pending', userId: 1 },
          { id: 2, title: 'Task 2', status: 'completed', userId: 2 }
        ],
        count: 2
      };

      const backendResponse = {
        tasks: [{ id: 1, title: 'Task 1', status: 'pending', userId: 1 }],
        count: 1
      };

      cache.get.mockReturnValue(JSON.stringify(cachedTasks));
      makeRequest.mockResolvedValue(backendResponse);

      await getAllTasks(mockReq, mockRes);

      expect(cache.set).not.toHaveBeenCalled();
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should handle malformed cached data', async () => {
      cache.get.mockReturnValue('invalid json');
      const backendResponse = { tasks: [], count: 0 };
      makeRequest.mockResolvedValue(backendResponse);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      expect(cache.set).toHaveBeenCalledWith('tasks', JSON.stringify(backendResponse));
    });

    it('should handle backend connection error', async () => {
      const connectionError = new Error('ECONNREFUSED');
      makeRequest.mockRejectedValue(connectionError);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'ECONNREFUSED',
        code: 'BACKEND_ERROR'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'GET',
        '/api/tasks',
        503,
        0,
        'GetTasksError',
        'ECONNREFUSED'
      );
    });
  });

  describe('getTaskById', () => {
    beforeEach(() => {
      mockReq.params = { id: '1' };
      mockReq.originalUrl = '/api/tasks/1';
    });

    it('should return cached task when available', async () => {
      const cachedTask = { id: 1, title: 'Task 1', status: 'pending', userId: 1 };
      cache.get.mockReturnValue(JSON.stringify(cachedTask));

      await getTaskById(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('tasks:1');
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.json).toHaveBeenCalledWith(cachedTask);
    });

    it('should fetch task from backend when not cached', async () => {
      const backendTask = { id: 1, title: 'Task 1', status: 'pending', userId: 1 };
      makeRequest.mockResolvedValue({ data: backendTask });
      cache.get.mockReturnValue(null);

      await getTaskById(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('tasks:1');
      expect(makeRequest).toHaveBeenCalledWith('/api/tasks/1');
      expect(cache.set).toHaveBeenCalledWith('tasks:1', JSON.stringify(backendTask));
      expect(mockRes.json).toHaveBeenCalledWith(backendTask);
    });

    it('should handle task not found error', async () => {
      const notFoundError = new Error('Task not found');
      notFoundError.statusCode = 404;
      makeRequest.mockRejectedValue(notFoundError);
      cache.get.mockReturnValue(null);

      await getTaskById(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(404);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Task not found',
        code: 'NOT_FOUND'
      });
    });
  });

  describe('createTask', () => {
    beforeEach(() => {
      mockReq.method = 'POST';
      mockReq.originalUrl = '/api/tasks';
      mockReq.body = {
        title: 'New Task',
        status: 'pending',
        userId: 1
      };
    });

    it('should create task successfully with valid data', async () => {
      const validation = { isValid: true, errors: [] };
      validateTask.mockReturnValue(validation);

      const createdTask = {
        id: 1,
        title: 'New Task',
        status: 'pending',
        userId: 1
      };
      makeRequest.mockResolvedValue({ data: createdTask });

      await createTask(mockReq, mockRes);

      expect(validateTask).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).toHaveBeenCalledWith('/api/tasks', {
        method: 'POST',
        body: mockReq.body
      });
      expect(cache.invalidate).toHaveBeenCalledWith('tasks');
      expect(mockRes.status).toHaveBeenCalledWith(201);
      expect(mockRes.json).toHaveBeenCalledWith(createdTask);
    });

    it('should return validation error for invalid data', async () => {
      const validation = { isValid: false, errors: ['Title is required', 'Invalid status'] };
      validateTask.mockReturnValue(validation);

      await createTask(mockReq, mockRes);

      expect(validateTask).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.status).toHaveBeenCalledWith(400);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Title is required, Invalid status'
      });
    });

    it('should handle backend validation error', async () => {
      const validation = { isValid: true, errors: [] };
      validateTask.mockReturnValue(validation);

      const backendError = new Error('User not found');
      backendError.statusCode = 400;
      makeRequest.mockRejectedValue(backendError);

      await createTask(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(400);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'User not found' });
    });
  });

  describe('updateTask', () => {
    beforeEach(() => {
      mockReq.method = 'PUT';
      mockReq.originalUrl = '/api/tasks/1';
      mockReq.params = { id: '1' };
      mockReq.body = {
        title: 'Updated Task',
        status: 'completed'
      };
    });

    it('should update task successfully with valid data', async () => {
      const validation = { isValid: true, errors: [] };
      validateTaskUpdate.mockReturnValue(validation);

      const updatedTask = {
        id: 1,
        title: 'Updated Task',
        status: 'completed',
        userId: 1
      };
      makeRequest.mockResolvedValue({ data: updatedTask });

      await updateTask(mockReq, mockRes);

      expect(validateTaskUpdate).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).toHaveBeenCalledWith('/api/tasks/1', {
        method: 'PUT',
        body: mockReq.body
      });
      expect(cache.invalidate).toHaveBeenCalledWith('tasks');
      expect(mockRes.json).toHaveBeenCalledWith(updatedTask);
    });

    it('should return validation error for invalid update data', async () => {
      const validation = { isValid: false, errors: ['Invalid status'] };
      validateTaskUpdate.mockReturnValue(validation);

      await updateTask(mockReq, mockRes);

      expect(validateTaskUpdate).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.status).toHaveBeenCalledWith(400);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Invalid status'
      });
    });

    it('should handle task not found during update', async () => {
      const validation = { isValid: true, errors: [] };
      validateTaskUpdate.mockReturnValue(validation);

      const notFoundError = new Error('Task not found');
      notFoundError.statusCode = 404;
      makeRequest.mockRejectedValue(notFoundError);

      await updateTask(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(404);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Task not found' });
    });
  });

  describe('deleteTask', () => {
    beforeEach(() => {
      mockReq.method = 'DELETE';
      mockReq.originalUrl = '/api/tasks/1';
      mockReq.params = { id: '1' };
    });

    it('should delete task successfully', async () => {
      makeRequest.mockResolvedValue({});

      await deleteTask(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks/1', {
        method: 'DELETE'
      });
      expect(cache.invalidate).toHaveBeenCalledWith('tasks');
      expect(mockRes.json).toHaveBeenCalledWith({ message: 'Task deleted successfully' });
    });

    it('should handle task not found during deletion', async () => {
      const notFoundError = new Error('Task not found');
      notFoundError.statusCode = 404;
      makeRequest.mockRejectedValue(notFoundError);

      await deleteTask(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(404);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Task not found' });
    });

    it('should handle backend server error during deletion', async () => {
      const serverError = new Error('Database error');
      serverError.statusCode = 500;
      makeRequest.mockRejectedValue(serverError);

      await deleteTask(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Database error' });
    });
  });

  describe('Error Handling Edge Cases', () => {
    it('should handle makeRequest throwing non-Error objects in getAllTasks', async () => {
      makeRequest.mockRejectedValue('String error');
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'String error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle cache errors gracefully', async () => {
      cache.get.mockImplementation(() => {
        throw new Error('Cache error');
      });

      const backendResponse = { tasks: [], count: 0 };
      makeRequest.mockResolvedValue(backendResponse);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should handle empty query parameters', async () => {
      mockReq.query = {};
      const backendResponse = { tasks: [], count: 0 };
      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllTasks(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      expect(cache.set).toHaveBeenCalledWith('tasks', JSON.stringify(backendResponse));
    });
  });
});
