const express = require('express');
const router = express.Router();
const cacheController = require('../controllers/cacheController');

router.get('/stats', cacheController.getCacheStats);

module.exports = router;
