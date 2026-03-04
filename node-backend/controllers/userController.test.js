const {
  getAllUsers,
  getUserById,
  createUser
} = require('./userController');
const { makeRequest } = require('../utils/httpClient');
const { validateUser } = require('../utils/validation');
const logger = require('../utils/logger');
const cache = require('../utils/cache');

// Mock dependencies
jest.mock('../utils/httpClient');
jest.mock('../utils/validation');
jest.mock('../utils/logger');
jest.mock('../utils/cache');

describe('User Controller', () => {
  let mockReq, mockRes;

  beforeEach(() => {
    jest.clearAllMocks();

    mockReq = {
      method: 'GET',
      originalUrl: '/api/users',
      params: {},
      body: {}
    };

    mockRes = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis()
    };
  });

  describe('getAllUsers', () => {
    it('should return cached users when available', async () => {
      const cachedUsers = {
        users: [
          { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' },
          { id: 2, name: 'Jane Smith', email: 'jane@example.com', role: 'designer' }
        ],
        count: 2
      };

      cache.get.mockReturnValue(JSON.stringify(cachedUsers));

      await getAllUsers(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('users');
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.json).toHaveBeenCalledWith(cachedUsers);
    });

    it('should fetch users from backend when not cached', async () => {
      const backendResponse = {
        users: [
          { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' }
        ],
        count: 1
      };

      makeRequest.mockResolvedValue(backendResponse);
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('users');
      expect(makeRequest).toHaveBeenCalledWith('/api/users');
      expect(cache.set).toHaveBeenCalledWith('users', JSON.stringify(backendResponse));
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should handle wrapped response from backend', async () => {
      const wrappedResponse = {
        data: {
          users: [{ id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' }],
          count: 1
        }
      };

      makeRequest.mockResolvedValue(wrappedResponse);
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      const expectedResponse = {
        users: [{ id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' }],
        count: 1
      };

      expect(mockRes.json).toHaveBeenCalledWith(expectedResponse);
    });

    it('should handle empty response from backend', async () => {
      makeRequest.mockResolvedValue({});
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      const expectedResponse = { users: [], count: 0 };
      expect(mockRes.json).toHaveBeenCalledWith(expectedResponse);
    });

    it('should handle backend connection error', async () => {
      const connectionError = new Error('ECONNREFUSED');
      makeRequest.mockRejectedValue(connectionError);
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(503);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'ECONNREFUSED',
        code: 'BACKEND_ERROR'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'GET',
        '/api/users',
        503,
        0,
        'GetUsersError',
        'ECONNREFUSED'
      );
    });

    it('should handle backend server error', async () => {
      const serverError = new Error('Internal Server Error');
      serverError.statusCode = 500;
      makeRequest.mockRejectedValue(serverError);
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle malformed cached data', async () => {
      cache.get.mockReturnValue('invalid json');

      await getAllUsers(mockReq, mockRes);

      // Should fallback to backend request
      expect(makeRequest).toHaveBeenCalledWith('/api/users');
    });
  });

  describe('getUserById', () => {
    beforeEach(() => {
      mockReq.params = { id: '1' };
      mockReq.originalUrl = '/api/users/1';
    });

    it('should return cached user when available', async () => {
      const cachedUser = { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' };
      cache.get.mockReturnValue(JSON.stringify(cachedUser));

      await getUserById(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('users:1');
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.json).toHaveBeenCalledWith(cachedUser);
    });

    it('should fetch user from backend when not cached', async () => {
      const backendUser = { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' };
      makeRequest.mockResolvedValue({ data: backendUser });
      cache.get.mockReturnValue(null);

      await getUserById(mockReq, mockRes);

      expect(cache.get).toHaveBeenCalledWith('users:1');
      expect(makeRequest).toHaveBeenCalledWith('/api/users/1');
      expect(cache.set).toHaveBeenCalledWith('users:1', JSON.stringify(backendUser));
      expect(mockRes.json).toHaveBeenCalledWith(backendUser);
    });

    it('should handle unwrapped response from backend', async () => {
      const backendUser = { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' };
      makeRequest.mockResolvedValue(backendUser);
      cache.get.mockReturnValue(null);

      await getUserById(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(backendUser);
    });

    it('should handle user not found error', async () => {
      const notFoundError = new Error('User not found');
      notFoundError.statusCode = 404;
      makeRequest.mockRejectedValue(notFoundError);
      cache.get.mockReturnValue(null);

      await getUserById(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(404);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'User not found',
        code: 'NOT_FOUND'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'GET',
        '/api/users/1',
        404,
        0,
        'GetUserByIdError',
        'User not found'
      );
    });

    it('should handle invalid user ID', async () => {
      mockReq.params = { id: 'invalid' };
      mockReq.originalUrl = '/api/users/invalid';

      const invalidIdError = new Error('Invalid user ID');
      invalidIdError.statusCode = 400;
      makeRequest.mockRejectedValue(invalidIdError);
      cache.get.mockReturnValue(null);

      await getUserById(mockReq, mockRes);

      expect(makeRequest).toHaveBeenCalledWith('/api/users/invalid');
      expect(mockRes.status).toHaveBeenCalledWith(400);
    });
  });

  describe('createUser', () => {
    beforeEach(() => {
      mockReq.method = 'POST';
      mockReq.originalUrl = '/api/users';
      mockReq.body = {
        name: 'John Doe',
        email: 'john@example.com',
        role: 'developer'
      };
    });

    it('should create user successfully with valid data', async () => {
      const validation = { isValid: true, errors: [] };
      validateUser.mockReturnValue(validation);

      const createdUser = {
        id: 1,
        name: 'John Doe',
        email: 'john@example.com',
        role: 'developer'
      };
      makeRequest.mockResolvedValue({ data: createdUser });

      await createUser(mockReq, mockRes);

      expect(validateUser).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).toHaveBeenCalledWith('/api/users', {
        method: 'POST',
        body: mockReq.body
      });
      expect(cache.invalidate).toHaveBeenCalledWith('users');
      expect(mockRes.status).toHaveBeenCalledWith(201);
      expect(mockRes.json).toHaveBeenCalledWith(createdUser);
    });

    it('should return validation error for invalid data', async () => {
      const validation = { isValid: false, errors: ['Name is required', 'Invalid email format'] };
      validateUser.mockReturnValue(validation);

      await createUser(mockReq, mockRes);

      expect(validateUser).toHaveBeenCalledWith(mockReq.body);
      expect(makeRequest).not.toHaveBeenCalled();
      expect(mockRes.status).toHaveBeenCalledWith(400);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Name is required, Invalid email format'
      });
      expect(logger.logError).toHaveBeenCalledWith(
        'POST',
        '/api/users',
        400,
        0,
        'ValidationError',
        'Name is required, Invalid email format'
      );
    });

    it('should handle backend validation error', async () => {
      const validation = { isValid: true, errors: [] };
      validateUser.mockReturnValue(validation);

      const backendError = new Error('Email already exists');
      backendError.statusCode = 400;
      makeRequest.mockRejectedValue(backendError);

      await createUser(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(400);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Email already exists' });
      expect(logger.logError).toHaveBeenCalledWith(
        'POST',
        '/api/users',
        400,
        0,
        'ValidationError',
        'Email already exists'
      );
    });

    it('should handle backend server error during creation', async () => {
      const validation = { isValid: true, errors: [] };
      validateUser.mockReturnValue(validation);

      const serverError = new Error('Database connection failed');
      serverError.statusCode = 500;
      makeRequest.mockRejectedValue(serverError);

      await createUser(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Database connection failed' });
      expect(logger.logError).toHaveBeenCalledWith(
        'POST',
        '/api/users',
        500,
        0,
        'CreateUserError',
        'Database connection failed'
      );
    });

    it('should handle unwrapped response from backend', async () => {
      const validation = { isValid: true, errors: [] };
      validateUser.mockReturnValue(validation);

      const createdUser = {
        id: 1,
        name: 'John Doe',
        email: 'john@example.com',
        role: 'developer'
      };
      makeRequest.mockResolvedValue(createdUser);

      await createUser(mockReq, mockRes);

      expect(mockRes.json).toHaveBeenCalledWith(createdUser);
    });

    it('should invalidate cache after successful creation', async () => {
      const validation = { isValid: true, errors: [] };
      validateUser.mockReturnValue(validation);

      const createdUser = { id: 1, name: 'John Doe' };
      makeRequest.mockResolvedValue({ data: createdUser });

      await createUser(mockReq, mockRes);

      expect(cache.invalidate).toHaveBeenCalledWith('users');
    });
  });

  describe('Error Handling Edge Cases', () => {
    it('should handle makeRequest throwing non-Error objects', async () => {
      makeRequest.mockRejectedValue('String error');
      cache.get.mockReturnValue(null);

      await getAllUsers(mockReq, mockRes);

      expect(mockRes.status).toHaveBeenCalledWith(500);
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'String error',
        code: 'BACKEND_ERROR'
      });
    });

    it('should handle cache.get throwing errors', async () => {
      cache.get.mockImplementation(() => {
        throw new Error('Cache error');
      });

      const backendResponse = { users: [], count: 0 };
      makeRequest.mockResolvedValue(backendResponse);

      await getAllUsers(mockReq, mockRes);

      // Should fallback to backend request
      expect(makeRequest).toHaveBeenCalledWith('/api/users');
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });

    it('should handle cache.set throwing errors', async () => {
      cache.get.mockReturnValue(null);
      cache.set.mockImplementation(() => {
        throw new Error('Cache set error');
      });

      const backendResponse = { users: [], count: 0 };
      makeRequest.mockResolvedValue(backendResponse);

      await getAllUsers(mockReq, mockRes);

      // Should still return the response despite cache error
      expect(mockRes.json).toHaveBeenCalledWith(backendResponse);
    });
  });
});
