const axios = require('axios');
const config = require('../config');
const logger = require('./logger');

/**
 * Helper function to make HTTP requests to Go backend
 * @param {string} path - API endpoint path
 * @param {object} options - Request options (method, headers, body)
 * @returns {Promise} - Resolves with response data or rejects with error
 */
async function makeRequest(path, options = {}) {
  const startTime = Date.now();
  const url = new URL(path, config.GO_BACKEND_URL).toString();
  const method = (options.method || 'GET').toUpperCase();

  try {
    const response = await axios({
      method,
      url,
      data: options.body,
      timeout: 5000,
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'node-backend/1.0.0',
        ...options.headers
      }
    });

    const duration = Date.now() - startTime;
    logger.logBackend(method, path, response.status, duration);
    return response.data;
  } catch (error) {
    const duration = Date.now() - startTime;
    const statusCode = error.response?.status || 0;
    const errMsg = error.response?.data?.error || error.response?.data?.message || error.message || String(error);
    logger.logError(`BACKEND ${method}`, path, statusCode, duration, error.code || error.name || 'RequestError', errMsg);

    if (error.response) {
      const err = new Error(error.message || `Request failed with status code ${error.response.status}`);
      err.statusCode = error.response.status;
      err.code = error.response.data?.code;
      throw err;
    }
    throw error;
  }
}

module.exports = { makeRequest };
