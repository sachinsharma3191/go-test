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
  if (!userData || typeof userData !== 'object') {
    errors.push('Name is required');
    errors.push('Email is required');
    errors.push('Role is required');
    return { isValid: errors.length === 0, errors };
  }

  // Check required fields
  if (!userData.name || typeof userData.name !== 'string' || userData.name.trim() === '') {
    errors.push('Name is required');
  }
  if (!userData.email || typeof userData.email !== 'string' || userData.email.trim() === '') {
    errors.push('Email is required');
  }
  if (!userData.role || typeof userData.role !== 'string' || userData.role.trim() === '') {
    errors.push('Role is required');
  }
  
  // Validate name length
  if (userData.name && typeof userData.name === 'string' && userData.name.length > 100) {
    errors.push('Name must be 100 characters or less');
  }
  
  // Validate email format (only if email is a string)
  if (userData.email && typeof userData.email === 'string' && userData.email.trim() !== '') {
    const email = userData.email.trim();
    if (email.includes('..')) {
      errors.push('Invalid email format');
    } else {
      const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      if (!emailRegex.test(email)) {
        errors.push('Invalid email format');
      }
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
  if (!taskData || typeof taskData !== 'object') {
    errors.push('Title is required');
    errors.push('Status is required');
    errors.push('User ID is required');
    return { isValid: errors.length === 0, errors };
  }

  // Check required fields
  if (!taskData.title || typeof taskData.title !== 'string' || taskData.title.trim() === '') {
    errors.push('Title is required');
  }
  if (!taskData.status || typeof taskData.status !== 'string' || taskData.status.trim() === '') {
    errors.push('Status is required');
  }
  if (taskData.userId === undefined || taskData.userId === null || taskData.userId === '') {
    errors.push('User ID is required');
  }
  
  // Validate status values
  if (taskData.status) {
    const validStatuses = ['pending', 'in-progress', 'completed'];
    if (!validStatuses.includes(taskData.status)) {
      errors.push('Status must be one of: pending, in-progress, completed');
    }
  }
  
  // Validate user ID is a positive integer
  if (taskData.userId !== undefined && taskData.userId !== null && taskData.userId !== '') {
    const num = Number(taskData.userId);
    if (isNaN(num) || num <= 0 || num % 1 !== 0) {
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
  if (!taskData || typeof taskData !== 'object') {
    return { isValid: true, errors: [] };
  }

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
  
  // Validate title if provided (cannot be empty when provided)
  if (taskData.title !== undefined && taskData.title !== null && typeof taskData.title === 'string' && taskData.title.trim() === '') {
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
