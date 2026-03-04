/**
 * Request logging middleware
 * Logs HTTP requests with method, path, status code, and duration
 */

const logger = require('../utils/logger');

const requestLogger = (req, res, next) => {
  const startTime = Date.now();
  
  // Store original end function
  const originalEnd = res.end;
  
  // Override end function to capture response details
  res.end = function(chunk, encoding) {
    // Calculate duration
    const duration = Date.now() - startTime;
    
    // Log the request using centralized logger
    logger.logRequest(req.method, req.originalUrl, res.statusCode, duration);
    
    // Call original end function
    originalEnd.call(this, chunk, encoding);
  };
  
  next();
};

module.exports = requestLogger;
