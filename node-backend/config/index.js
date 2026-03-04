module.exports = {
  get PORT() { return process.env.PORT || 3000; },
  get GO_BACKEND_URL() { return process.env.GO_BACKEND_URL || 'http://localhost:8080'; },
  get NODE_ENV() { return process.env.NODE_ENV || 'development'; }
};
