const request = require('supertest');
const app = require('../app');
const cache = require('../utils/cache');

// Mock the Go backend HTTP client for integration tests
jest.mock('../utils/httpClient', () => ({
  makeRequest: jest.fn()
}));

const { makeRequest } = require('../utils/httpClient');

describe('Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    if (cache.clear) cache.clear();
  });

  describe('Health Endpoint Integration', () => {
    it('should return healthy status when backend is available', async () => {
      makeRequest.mockResolvedValue({
        status: 'healthy',
        message: 'All systems operational'
      });

      const response = await request(app)
        .get('/health')
        .expect(200);

      expect(response.body).toEqual({
        status: 'ok',
        message: 'Node.js backend is running',
        goBackend: {
          status: 'healthy',
          message: 'All systems operational'
        }
      });

      expect(makeRequest).toHaveBeenCalledWith('/health');
    });

    it('should return error status when backend is unavailable', async () => {
      makeRequest.mockRejectedValue(new Error('ECONNREFUSED'));

      const response = await request(app)
        .get('/health')
        .expect(503);

      expect(response.body).toEqual({
        status: 'error',
        message: 'Node.js backend is running but Go backend is unavailable',
        error: 'ECONNREFUSED'
      });
    });
  });

  describe('User Management Integration', () => {
    describe('GET /api/users', () => {
      it('should return users list successfully', async () => {
        const mockUsers = {
          users: [
            { id: 1, name: 'John Doe', email: 'john@example.com', role: 'developer' },
            { id: 2, name: 'Jane Smith', email: 'jane@example.com', role: 'designer' }
          ],
          count: 2
        };

        makeRequest.mockResolvedValue(mockUsers);

        const response = await request(app)
          .get('/api/users')
          .expect(200);

        expect(response.body).toEqual(mockUsers);
        expect(makeRequest).toHaveBeenCalledWith('/api/users');
      });

      it('should handle backend errors gracefully', async () => {
        makeRequest.mockRejectedValue(new Error('Backend unavailable'));

        const response = await request(app)
          .get('/api/users')
          .expect(503);

        expect(response.body).toEqual({
          error: 'Backend unavailable',
          code: 'BACKEND_ERROR'
        });
      });

      it('should handle empty users list', async () => {
        makeRequest.mockResolvedValue({ users: [], count: 0 });

        const response = await request(app)
          .get('/api/users')
          .expect(200);

        expect(response.body).toEqual({ users: [], count: 0 });
      });
    });

    describe('GET /api/users/:id', () => {
      it('should return specific user', async () => {
        const mockUser = {
          id: 1,
          name: 'John Doe',
          email: 'john@example.com',
          role: 'developer'
        };

        makeRequest.mockResolvedValue({ data: mockUser });

        const response = await request(app)
          .get('/api/users/1')
          .expect(200);

        expect(response.body).toEqual(mockUser);
        expect(makeRequest).toHaveBeenCalledWith('/api/users/1');
      });

      it('should handle user not found', async () => {
        const notFoundError = new Error('User not found');
        notFoundError.statusCode = 404;
        makeRequest.mockRejectedValue(notFoundError);

        const response = await request(app)
          .get('/api/users/999')
          .expect(404);

        expect(response.body).toEqual({
          error: 'User not found',
          code: 'NOT_FOUND'
        });
      });

      it('should handle invalid user ID', async () => {
        const invalidIdError = new Error('Invalid user ID');
        invalidIdError.statusCode = 400;
        makeRequest.mockRejectedValue(invalidIdError);

        const response = await request(app)
          .get('/api/users/invalid')
          .expect(400);

        expect(response.body).toEqual({
          error: 'Invalid user ID',
          code: 'NOT_FOUND'
        });
      });
    });

    describe('POST /api/users', () => {
      it('should create new user successfully', async () => {
        const newUser = {
          name: 'Alice Johnson',
          email: 'alice@example.com',
          role: 'manager'
        };

        const createdUser = {
          id: 3,
          ...newUser
        };

        makeRequest.mockResolvedValue({ data: createdUser });

        const response = await request(app)
          .post('/api/users')
          .send(newUser)
          .expect(201);

        expect(response.body).toEqual(createdUser);
        expect(makeRequest).toHaveBeenCalledWith('/api/users', {
          method: 'POST',
          body: newUser
        });
      });

      it('should handle validation errors', async () => {
        const invalidUser = {
          name: '', // Empty name
          email: 'invalid-email', // Invalid email
          role: '' // Empty role
        };

        const response = await request(app)
          .post('/api/users')
          .send(invalidUser)
          .expect(400);

        expect(response.body.error).toContain('Name is required');
        expect(response.body.error).toContain('Role is required');
        expect(response.body.error).toContain('Invalid email format');
        expect(makeRequest).not.toHaveBeenCalled();
      });

      it('should handle backend validation errors', async () => {
        const userWithDuplicateEmail = {
          name: 'John Doe',
          email: 'john@example.com', // Duplicate email
          role: 'developer'
        };

        const duplicateError = new Error('Email already exists');
        duplicateError.statusCode = 400;
        makeRequest.mockRejectedValue(duplicateError);

        const response = await request(app)
          .post('/api/users')
          .send(userWithDuplicateEmail)
          .expect(400);

        expect(response.body).toEqual({
          error: 'Email already exists'
        });
      });
    });
  });

  describe('Task Management Integration', () => {
    describe('GET /api/tasks', () => {
      it('should return tasks list successfully', async () => {
        const mockTasks = {
          tasks: [
            { id: 1, title: 'Task 1', status: 'pending', userId: 1 },
            { id: 2, title: 'Task 2', status: 'completed', userId: 2 }
          ],
          count: 2
        };

        makeRequest.mockResolvedValue(mockTasks);

        const response = await request(app)
          .get('/api/tasks')
          .expect(200);

        expect(response.body).toEqual(mockTasks);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks');
      });

      it('should handle status filter', async () => {
        const mockFilteredTasks = {
          tasks: [
            { id: 1, title: 'Task 1', status: 'pending', userId: 1 }
          ],
          count: 1
        };

        makeRequest.mockResolvedValue(mockFilteredTasks);

        const response = await request(app)
          .get('/api/tasks?status=pending')
          .expect(200);

        expect(response.body).toEqual(mockFilteredTasks);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks?status=pending');
      });

      it('should handle userId filter', async () => {
        const mockUserTasks = {
          tasks: [
            { id: 1, title: 'Task 1', status: 'pending', userId: 1 },
            { id: 2, title: 'Task 2', status: 'completed', userId: 1 }
          ],
          count: 2
        };

        makeRequest.mockResolvedValue(mockUserTasks);

        const response = await request(app)
          .get('/api/tasks?userId=1')
          .expect(200);

        expect(response.body).toEqual(mockUserTasks);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks?userId=1');
      });

      it('should handle multiple filters', async () => {
        const mockFilteredTasks = {
          tasks: [
            { id: 1, title: 'Task 1', status: 'pending', userId: 1 }
          ],
          count: 1
        };

        makeRequest.mockResolvedValue(mockFilteredTasks);

        const response = await request(app)
          .get('/api/tasks?status=pending&userId=1')
          .expect(200);

        expect(response.body).toEqual(mockFilteredTasks);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks?status=pending&userId=1');
      });
    });

    describe('POST /api/tasks', () => {
      it('should create new task successfully', async () => {
        const newTask = {
          title: 'Implement new feature',
          status: 'pending',
          userId: 1
        };

        const createdTask = {
          id: 3,
          ...newTask
        };

        makeRequest.mockResolvedValue({ data: createdTask });

        const response = await request(app)
          .post('/api/tasks')
          .send(newTask)
          .expect(201);

        expect(response.body).toEqual(createdTask);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks', {
          method: 'POST',
          body: newTask
        });
      });

      it('should handle task validation errors', async () => {
        const invalidTask = {
          title: '', // Empty title
          status: 'invalid-status', // Invalid status
          userId: -1 // Invalid user ID
        };

        const response = await request(app)
          .post('/api/tasks')
          .send(invalidTask)
          .expect(400);

        expect(response.body).toEqual({
          error: 'Title is required, Status must be one of: pending, in-progress, completed, User ID must be a positive number'
        });
        expect(makeRequest).not.toHaveBeenCalled();
      });
    });

    describe('PUT /api/tasks/:id', () => {
      it('should update task successfully', async () => {
        const taskUpdate = {
          title: 'Updated task title',
          status: 'completed'
        };

        const updatedTask = {
          id: 1,
          title: 'Updated task title',
          status: 'completed',
          userId: 1
        };

        makeRequest.mockResolvedValue({ data: updatedTask });

        const response = await request(app)
          .put('/api/tasks/1')
          .send(taskUpdate)
          .expect(200);

        expect(response.body).toEqual(updatedTask);
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks/1', {
          method: 'PUT',
          body: taskUpdate
        });
      });

      it('should handle task not found during update', async () => {
        const taskUpdate = {
          title: 'Updated task',
          status: 'completed'
        };

        const notFoundError = new Error('Task not found');
        notFoundError.statusCode = 404;
        makeRequest.mockRejectedValue(notFoundError);

        const response = await request(app)
          .put('/api/tasks/999')
          .send(taskUpdate)
          .expect(404);

        expect(response.body).toEqual({
          error: 'Task not found'
        });
      });
    });

    describe('DELETE /api/tasks/:id', () => {
      it('should delete task successfully', async () => {
        makeRequest.mockResolvedValue({});

        const response = await request(app)
          .delete('/api/tasks/1')
          .expect(200);

        expect(response.body).toEqual({
          message: 'Task deleted successfully'
        });
        expect(makeRequest).toHaveBeenCalledWith('/api/tasks/1', {
          method: 'DELETE'
        });
      });

      it('should handle task not found during deletion', async () => {
        const notFoundError = new Error('Task not found');
        notFoundError.statusCode = 404;
        makeRequest.mockRejectedValue(notFoundError);

        const response = await request(app)
          .delete('/api/tasks/999')
          .expect(404);

        expect(response.body).toEqual({
          error: 'Task not found'
        });
      });
    });
  });

  describe('Statistics Integration', () => {
    describe('GET /api/stats', () => {
      it('should return system statistics', async () => {
        const mockStats = {
          users: { total: 10 },
          tasks: {
            total: 25,
            pending: 8,
            inProgress: 12,
            completed: 5
          }
        };

        makeRequest.mockResolvedValue({ data: mockStats });

        const response = await request(app)
          .get('/api/stats')
          .expect(200);

        expect(response.body).toEqual(mockStats);
        expect(makeRequest).toHaveBeenCalledWith('/api/stats');
      });

      it('should handle stats endpoint errors', async () => {
        makeRequest.mockRejectedValue(new Error('Stats service unavailable'));

        const response = await request(app)
          .get('/api/stats')
          .expect(503);

        expect(response.body).toEqual({
          error: 'Stats service unavailable',
          code: 'BACKEND_ERROR'
        });
      });
    });

    describe('GET /api/stats/cache', () => {
      it('should return cache statistics', async () => {
        const mockCacheStats = {
          hits: 150,
          misses: 25,
          evictions: 0,
          totalEntries: 10
        };

        makeRequest.mockResolvedValue({ data: mockCacheStats });

        const response = await request(app)
          .get('/api/stats/cache')
          .expect(200);

        expect(response.body).toEqual(mockCacheStats);
        expect(makeRequest).toHaveBeenCalledWith('/api/stats/cache');
      });

      it('should handle cache stats errors', async () => {
        makeRequest.mockRejectedValue(new Error('Cache service unavailable'));

        const response = await request(app)
          .get('/api/stats/cache')
          .expect(503);

        expect(response.body).toEqual({
          error: 'Cache service unavailable',
          code: 'BACKEND_ERROR'
        });
      });
    });
  });

  describe('Error Handling Integration', () => {
    it('should handle malformed JSON requests', async () => {
      const response = await request(app)
        .post('/api/users')
        .set('Content-Type', 'application/json')
        .send('invalid-json')
        .expect(400);

      expect(response.body).toHaveProperty('error');
    });

    it('should handle missing content-type header', async () => {
      const response = await request(app)
        .post('/api/users')
        .send('{ "name": "test" }')
        .expect(400);

      expect(response.body).toHaveProperty('error');
    });

    it('should handle empty request body', async () => {
      const response = await request(app)
        .post('/api/users')
        .send({})
        .expect(400);

      expect(response.body).toEqual({
        error: 'Name is required, Email is required, Role is required'
      });
    });

    it('should handle very large request bodies', async () => {
      const largeData = {
        name: 'A'.repeat(1000),
        email: 'test@example.com',
        role: 'developer'
      };

      const response = await request(app)
        .post('/api/users')
        .send(largeData)
        .expect(400);

      expect(response.body).toHaveProperty('error');
    });
  });

  describe('CORS Integration', () => {
    it('should include CORS headers', async () => {
      makeRequest.mockResolvedValue({ status: 'healthy' });

      const response = await request(app)
        .get('/health')
        .set('Origin', 'http://localhost:3000')
        .expect(200);

      // CORS adds allow-origin for requests with Origin; methods/headers on preflight
      expect(response.headers['access-control-allow-origin']).toBeDefined();
    });

    it('should handle OPTIONS requests', async () => {
      const response = await request(app)
        .options('/api/users')
        .set('Origin', 'http://localhost:3000')
        .expect(204);

      expect(response.headers['access-control-allow-origin']).toBeDefined();
      expect(response.headers['access-control-allow-methods']).toBeDefined();
    });
  });

  describe('Request Logging Integration', () => {
    it('should log requests appropriately', async () => {
      makeRequest.mockResolvedValue({ users: [], count: 0 });

      const consoleSpy = jest.spyOn(console, 'log').mockImplementation();

      await request(app)
        .get('/api/users')
        .expect(200);

      // Logger format: timestamp GET /api/users 200 <duration> (no ms suffix)
      expect(consoleSpy).toHaveBeenCalledWith(
        expect.stringMatching(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} GET \/api\/users 200 \d+$/)
      );

      consoleSpy.mockRestore();
    });

    it('should log errors appropriately', async () => {
      makeRequest.mockRejectedValue(new Error('Backend unavailable'));

      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

      await request(app)
        .get('/api/users')
        .expect(503);

      // Logger format: timestamp ERROR method path statusCode duration - method= method error=errorName message='...'
      expect(consoleSpy).toHaveBeenCalled();
      const logMsg = consoleSpy.mock.calls[0][0];
      expect(logMsg).toContain('ERROR');
      expect(logMsg).toContain('/api/users');
      expect(logMsg).toContain('503');
      expect(logMsg).toContain('GetUsersError');
      expect(logMsg).toContain('Backend');

      consoleSpy.mockRestore();
    });
  });

  describe('Performance Integration', () => {
    it('should handle concurrent requests', async () => {
      cache.clear();
      // Delay mock so all 10 requests pass cache check before any resolve
      makeRequest.mockImplementation(() => new Promise(resolve => setTimeout(() => resolve({ users: [], count: 0 }), 20)));

      const promises = [];
      const concurrentRequests = 10;

      for (let i = 0; i < concurrentRequests; i++) {
        promises.push(
          request(app)
            .get('/api/users')
            .expect(200)
        );
      }

      const responses = await Promise.all(promises);

      // All requests should succeed
      responses.forEach(response => {
        expect(response.body).toEqual({ users: [], count: 0 });
      });

      // With delayed resolution, all requests should miss cache and call backend
      expect(makeRequest).toHaveBeenCalledTimes(concurrentRequests);
    });

    it('should complete requests within reasonable time', async () => {
      makeRequest.mockResolvedValue({ users: [], count: 0 });

      const startTime = Date.now();

      await request(app)
        .get('/api/users')
        .expect(200);

      const endTime = Date.now();
      const duration = endTime - startTime;

      // Should complete within 100ms (adjust based on your performance requirements)
      expect(duration).toBeLessThan(100);
    });
  });

  describe('Content Type Handling', () => {
    it('should accept application/json content type', async () => {
      const newUser = {
        name: 'John Doe',
        email: 'john@example.com',
        role: 'developer'
      };
      const createdUser = { id: 1, ...newUser };
      makeRequest.mockResolvedValue({ data: createdUser });

      const response = await request(app)
        .post('/api/users')
        .set('Content-Type', 'application/json')
        .send(newUser)
        .expect(201);

      expect(response.body).toHaveProperty('id');
      expect(response.body).toHaveProperty('name', 'John Doe');
    });

    it('should reject invalid content types', async () => {
      const response = await request(app)
        .post('/api/users')
        .set('Content-Type', 'text/plain')
        .send('plain text data')
        .expect(400);

      expect(response.body).toHaveProperty('error');
    });
  });

  describe('URL Parameter Handling', () => {
    it('should handle special characters in URLs', async () => {
      makeRequest.mockResolvedValue({ data: { id: 1, name: 'John Doe' } });

      const response = await request(app)
        .get('/api/users/john.doe@example.com')
        .expect(200);

      expect(response.body).toHaveProperty('name', 'John Doe');
      expect(makeRequest).toHaveBeenCalledWith('/api/users/john.doe@example.com');
    });

    it('should handle URL encoded parameters', async () => {
      makeRequest.mockResolvedValue({
        tasks: [{ id: 1, title: 'Task with spaces', status: 'pending' }],
        count: 1
      });

      const response = await request(app)
        .get('/api/tasks?title=Task%20with%20spaces')
        .expect(200);

      expect(response.body).toHaveProperty('tasks');
      // URLSearchParams may use + for spaces; both are valid
      const calledUrl = makeRequest.mock.calls[0][0];
      expect(calledUrl).toMatch(/\/api\/tasks\?title=Task/);
      expect(calledUrl).toMatch(/with%20spaces|with\+spaces/);
    });
  });
});
