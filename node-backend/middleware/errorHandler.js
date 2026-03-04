/**
 * Global error handling middleware
 */
const logger = require('../utils/logger');

const errorHandler = (err, req, res, next) => {
  logger.logError(req.method, req.originalUrl, err.statusCode || 500, 0, err.message);
  console.error('Stack trace:', err.stack);
  
  res.status(err.statusCode || 500).json({ 
    error: err.message || 'Something went wrong!' 
  });
};

/**
 * 404 Not Found handler
 */
const notFoundHandler = (req, res) => {
  logger.logWarning(req.method, req.originalUrl, 404, 0, 'Route not found');
  
  res.status(404).json({ error: 'Route not found' });
};

module.exports = {
  errorHandler,
  notFoundHandler
};
