const express = require('express');
const cors = require('cors');
const config = require('./config');
const routes = require('./routes');
const { errorHandler, notFoundHandler } = require('./middleware/errorHandler');
const requestLogger = require('./middleware/logger');
const logger = require('./utils/logger');

const app = express();

// Middleware
app.use(cors());
app.use(express.json());
app.use(requestLogger); // Add logging middleware

// Routes
app.use(routes);

// Error handling
app.use(errorHandler);
app.use(notFoundHandler);

// Start server only when not in test (supertest uses app directly)
if (process.env.NODE_ENV !== 'test') {
  app.listen(config.PORT, () => {
    logger.logStartup(config.PORT, config.GO_BACKEND_URL);
  });
}

module.exports = app;
