/**
 * Centralized logging utility for Node.js backend
 * Matches the format of Go backend logging
 */

class Logger {
  constructor() {
    this.logLevel = process.env.LOG_LEVEL || 'info';
  }

  /**
   * Format log message with timestamp and structured data
   */
  formatMessage(level, method, path, statusCode, duration, message = '') {
    const timestamp = new Date().toISOString().replace('T', ' ').substring(0, 19);
    const baseMessage = `${timestamp} ${method} ${path} ${statusCode} ${duration}`;
    return message ? `${baseMessage} - ${message}` : baseMessage;
  }

  /**
   * Log successful requests
   */
  logRequest(method, path, statusCode, duration) {
    const message = this.formatMessage('INFO', method, path, statusCode, duration);
    console.log(message);
  }

  /**
   * Log backend requests
   */
  logBackend(method, path, statusCode, duration) {
    const message = this.formatMessage('INFO', `BACKEND ${method}`, path, statusCode, duration);
    console.log(message);
  }

  /**
   * Log errors with enhanced details
   */
  logError(method, path, statusCode, duration, errorName, errorMessage) {
    const timestamp = new Date().toISOString().replace('T', ' ').substring(0, 19);
    const message = `${timestamp} ERROR ${method} ${path} ${statusCode} ${duration} - method=${method} error=${errorName} message='${errorMessage}'`;
    console.error(message);
  }

  /**
   * Log warnings
   */
  logWarning(method, path, statusCode, duration, message) {
    const logMessage = this.formatMessage('WARNING', method, path, statusCode, duration, message);
    console.warn(logMessage);
  }

  /**
   * Log cache statistics
   */
  logCacheStats(stats) {
    const timestamp = new Date().toISOString().replace('T', ' ').substring(0, 19);
    const s = stats || {};
    const message = `${timestamp} CACHE_STATS hits=${s.hits} misses=${s.misses} evictions=${s.evictions} entries=${s.totalEntries}`;
    console.log(message);
  }

  /**
   * Log server startup
   */
  logStartup(port, goBackendUrl) {
    const timestamp = new Date().toISOString().replace('T', ' ').substring(0, 19);
    console.log(`${timestamp} Node.js backend server starting on http://localhost:${port}`);
    console.log(`${timestamp} Connecting to Go backend at ${goBackendUrl}`);
    console.log(`${timestamp} Request logging enabled`);
    console.log(`${timestamp} Health check: http://localhost:${port}/health`);
  }

  /**
   * Log data persistence events
   */
  logDataPersistence(action, file, error = null) {
    const timestamp = new Date().toISOString().replace('T', ' ').substring(0, 19);
    if (error) {
      console.error(`${timestamp} DATA_PERSISTENCE ${action} ${file} ERROR - ${error}`);
    } else {
      console.log(`${timestamp} DATA_PERSISTENCE ${action} ${file} SUCCESS`);
    }
  }
}

module.exports = new Logger();
