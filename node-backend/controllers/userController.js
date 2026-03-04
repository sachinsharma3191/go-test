const { makeRequest } = require('../utils/httpClient');
const { validateUser } = require('../utils/validation');
const logger = require('../utils/logger');
const cache = require('../utils/cache');

const CACHE_KEY_USERS = 'users';

/**
 * Get all users
 */
const getAllUsers = async (req, res) => {
  try {
    let cached;
    try {
      cached = cache.get(CACHE_KEY_USERS);
    } catch {
      cached = null;
    }
    if (cached) {
      try {
        return res.json(JSON.parse(cached));
      } catch {
        // Malformed cached data, fallback to backend
      }
    }

    const response = await makeRequest('/api/users');
    
    // Handle empty or wrapped response: return shape frontend expects { users, count }
    const users = response?.users ?? response?.data?.users ?? [];
    const count = response?.count ?? response?.data?.count ?? users.length;
    const result = { users, count };
    try {
      cache.set(CACHE_KEY_USERS, JSON.stringify(result));
    } catch {
      // Ignore cache set errors, still return response
    }
    return res.json(result);
  } catch (error) {
    const errMsg = error?.message ?? String(error);
    const statusCode = error?.statusCode || (errMsg.includes('ECONNREFUSED') || errMsg.includes('unavailable') ? 503 : 500);
    logger.logError(req.method, req.originalUrl, statusCode, 0, 'GetUsersError', errMsg);
    res.status(statusCode).json({
      error: errMsg || 'Failed to fetch users',
      code: error?.code || 'BACKEND_ERROR'
    });
  }
};

/**
 * Get user by ID
 */
const getUserById = async (req, res) => {
  try {
    const cacheKey = `users:${req.params.id}`;
    const cached = cache.get(cacheKey);
    if (cached) {
      return res.json(JSON.parse(cached));
    }

    const raw = await makeRequest(`/api/users/${req.params.id}`);
    const user = raw?.data ?? raw;
    cache.set(cacheKey, JSON.stringify(user));
    res.json(user);
  } catch (error) {
    const errMsg = error?.message ?? String(error);
    const statusCode = error?.statusCode || (errMsg.includes('not found') ? 404 : 500);
    logger.logError(req.method, req.originalUrl, statusCode, 0, 'GetUserByIdError', errMsg);
    res.status(statusCode).json({
      error: errMsg || 'User not found',
      code: error?.code || 'NOT_FOUND'
    });
  }
};

/**
 * Create a new user
 */
const createUser = async (req, res) => {
  try {
    // Validate input
    const validation = validateUser(req.body);
    if (!validation.isValid) {
      logger.logError(req.method, req.originalUrl, 400, 0, 'ValidationError', validation.errors.join(', '));
      return res.status(400).json({ 
        error: validation.errors.join(', ') 
      });
    }
    
    const raw = await makeRequest('/api/users', {
      method: 'POST',
      body: req.body
    });
    const user = raw?.data ?? raw;
    cache.invalidate('users');
    res.status(201).json(user);
  } catch (error) {
    const errMsg = error?.message ?? String(error);
    const statusCode = error?.statusCode || (errMsg.includes('400') ? 400 : 500);
    const errorName = statusCode === 400 ? 'ValidationError' : 'CreateUserError';
    logger.logError(req.method, req.originalUrl, statusCode, 0, errorName, errMsg);
    res.status(statusCode).json({ error: errMsg });
  }
};

module.exports = {
  getAllUsers,
  getUserById,
  createUser
};
