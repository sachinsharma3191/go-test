/**
 * Validation utilities for user and task data
 */

/**
 * Validate user creation data
 * @param {Object} userData - User data to validate
 * @returns {Object} - { isValid: boolean, errors: string[] }
 */
const validateUser = (userData) => {
  const errors = [];
  
  // Check required fields
  if (!userData.name || userData.name.trim() === '') {
    errors.push('Name is required');
  }
  
  if (!userData.email || userData.email.trim() === '') {
    errors.push('Email is required');
  }
  
  if (!userData.role || userData.role.trim() === '') {
    errors.push('Role is required');
  }
  
  // Validate name length
  if (userData.name && userData.name.length > 100) {
    errors.push('Name must be 100 characters or less');
  }
  
  // Validate email format
  if (userData.email) {
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(userData.email)) {
      errors.push('Invalid email format');
    }
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

/**
 * Validate task creation data
 * @param {Object} taskData - Task data to validate
 * @returns {Object} - { isValid: boolean, errors: string[] }
 */
const validateTask = (taskData) => {
  const errors = [];
  
  // Check required fields
  if (!taskData.title || taskData.title.trim() === '') {
    errors.push('Title is required');
  }
  
  if (!taskData.status || taskData.status.trim() === '') {
    errors.push('Status is required');
  }
  
  if (!taskData.userId || taskData.userId === '') {
    errors.push('User ID is required');
  }
  
  // Validate status values
  if (taskData.status) {
    const validStatuses = ['pending', 'in-progress', 'completed'];
    if (!validStatuses.includes(taskData.status)) {
      errors.push('Status must be one of: pending, in-progress, completed');
    }
  }
  
  // Validate user ID is a number
  if (taskData.userId) {
    const userId = parseInt(taskData.userId);
    if (isNaN(userId) || userId <= 0) {
      errors.push('User ID must be a positive number');
    }
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

/**
 * Validate task update data
 * @param {Object} taskData - Task data to validate
 * @returns {Object} - { isValid: boolean, errors: string[] }
 */
const validateTaskUpdate = (taskData) => {
  const errors = [];
  
  // Validate status if provided
  if (taskData.status) {
    const validStatuses = ['pending', 'in-progress', 'completed'];
    if (!validStatuses.includes(taskData.status)) {
      errors.push('Status must be one of: pending, in-progress, completed');
    }
  }
  
  // Validate user ID if provided
  if (taskData.userId) {
    const userId = parseInt(taskData.userId);
    if (isNaN(userId) || userId <= 0) {
      errors.push('User ID must be a positive number');
    }
  }
  
  // Validate title if provided
  if (taskData.title && taskData.title.trim() === '') {
    errors.push('Title cannot be empty');
  }
  
  return {
    isValid: errors.length === 0,
    errors
  };
};

module.exports = {
  validateUser,
  validateTask,
  validateTaskUpdate
};
