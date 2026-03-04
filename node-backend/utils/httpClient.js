const http = require('http');
const config = require('../config');
const logger = require('./logger');

/**
 * Helper function to make HTTP requests to Go backend
 * @param {string} path - API endpoint path
 * @param {object} options - Request options (method, headers, body)
 * @returns {Promise} - Resolves with response data or rejects with error
 */
function makeRequest(path, options = {}) {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    const url = new URL(path, config.GO_BACKEND_URL);
    const requestOptions = {
      hostname: url.hostname,
      port: url.port || 8080,
      path: url.pathname + url.search,
      method: options.method || 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      }
    };

    const req = http.request(requestOptions, (res) => {
      let data = '';
      res.on('data', (chunk) => {
        data += chunk;
      });
      res.on('end', () => {
        const duration = Date.now() - startTime;
        
        if (res.statusCode >= 200 && res.statusCode < 300) {
          logger.logBackend(requestOptions.method, path, res.statusCode, duration);
          try {
            resolve(JSON.parse(data));
          } catch (e) {
            resolve(data);
          }
        } else {
          logger.logError(`BACKEND ${requestOptions.method}`, path, res.statusCode, duration, data);
          const error = new Error(data || `Request failed with status ${res.statusCode}`);
          error.statusCode = res.statusCode;
          error.responseData = data;
          try {
            const errorData = JSON.parse(data);
            error.code = errorData.code;
            error.message = errorData.error || errorData.message || error.message;
          } catch (e) {
            // Keep error.message from data if not JSON
          }
          reject(error);
        }
      });
    });

    req.on('error', (error) => {
      const duration = Date.now() - startTime;
      logger.logError(`BACKEND ${requestOptions.method}`, path, 0, duration, error.message);
      reject(error);
    });

    if (options.body) {
      req.write(JSON.stringify(options.body));
    }

    req.end();
  });
}

module.exports = { makeRequest };
