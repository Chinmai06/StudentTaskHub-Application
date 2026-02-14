import React, { useState } from 'react';
import './CreateTask.css';

const CreateTask = () => {
  const [formData, setFormData] = useState({
    title: '',
    description1: '',
    description2: '',
    dueDate: '',
    priority: '',
    category: 'academic',
    attachments1: [],
    attachments2: []
  });

  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const handlePriorityChange = (e) => {
    setFormData(prev => ({
      ...prev,
      priority: e.target.checked ? 'high' : ''
    }));
  };

  const handleFileUpload1 = (e) => {
    const files = Array.from(e.target.files);
    setFormData(prev => ({
      ...prev,
      attachments1: [...prev.attachments1, ...files]
    }));
  };

  const handleFileUpload2 = (e) => {
    const files = Array.from(e.target.files);
    setFormData(prev => ({
      ...prev,
      attachments2: [...prev.attachments2, ...files]
    }));
  };

  const removeAttachment1 = (index) => {
    setFormData(prev => ({
      ...prev,
      attachments1: prev.attachments1.filter((_, i) => i !== index)
    }));
  };

  const removeAttachment2 = (index) => {
    setFormData(prev => ({
      ...prev,
      attachments2: prev.attachments2.filter((_, i) => i !== index)
    }));
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.title.trim()) {
      newErrors.title = 'Task title is required';
    }
    
    if (!formData.dueDate) {
      newErrors.dueDate = 'Due date is required';
    } else {
      const selectedDate = new Date(formData.dueDate);
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      if (selectedDate < today) {
        newErrors.dueDate = 'Due date cannot be in the past';
      }
    }
    
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const newErrors = validateForm();
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsSubmitting(true);
    
    try {
      const response = await fetch('/api/tasks', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...formData,
          createdAt: new Date().toISOString()
        })
      });

      if (!response.ok) {
        throw new Error('Failed to create task');
      }

      const data = await response.json();
      console.log('Task created:', data);
      
      setSubmitSuccess(true);
      
      setTimeout(() => {
        setFormData({
          title: '',
          description1: '',
          description2: '',
          dueDate: '',
          priority: '',
          category: 'academic',
          attachments1: [],
          attachments2: []
        });
        setSubmitSuccess(false);
      }, 2000);
      
    } catch (error) {
      console.error('Error:', error);
      setErrors({ submit: 'Failed to create task. Please try again.' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleReset = () => {
    setFormData({
      title: '',
      description1: '',
      description2: '',
      dueDate: '',
      priority: '',
      category: 'academic',
      attachments1: [],
      attachments2: []
    });
    setErrors({});
    setSubmitSuccess(false);
  };

  return (
    <div className="create-task-container">
      <div className="create-task-card">
        <div className="card-header">
          <h1>Create New Task</h1>
          <p className="subtitle">Stay organized, achieve your goals</p>
        </div>

        <form onSubmit={handleSubmit} className="task-form">
          {/* Task Title */}
          <div className="form-group">
            <label htmlFor="title">
              Task Title <span className="required">*</span>
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              placeholder="e.g., Complete Math Assignment"
              className={errors.title ? 'error' : ''}
            />
            {errors.title && <span className="error-message">{errors.title}</span>}
          </div>

          {/* First Description with Buttons */}
          <div className="form-group">
            <label htmlFor="description1">Description</label>
            <div className="description-wrapper">
              <textarea
                id="description1"
                name="description1"
                value={formData.description1}
                onChange={handleChange}
                placeholder="Add any additional details about this task..."
                rows="4"
              />
              <div className="description-actions">
                <input
                  type="file"
                  id="imageInput1"
                  accept="image/*"
                  multiple
                  onChange={handleFileUpload1}
                  style={{ display: 'none' }}
                />
                <label htmlFor="imageInput1" className="action-btn" title="Add image">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                    <circle cx="8.5" cy="8.5" r="1.5"></circle>
                    <polyline points="21 15 16 10 5 21"></polyline>
                  </svg>
                </label>

                <input
                  type="file"
                  id="fileInput1"
                  multiple
                  onChange={handleFileUpload1}
                  style={{ display: 'none' }}
                />
                <label htmlFor="fileInput1" className="action-btn" title="Add file">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                    <polyline points="13 2 13 9 20 9"></polyline>
                  </svg>
                </label>

                <input
                  type="file"
                  id="folderInput1"
                  webkitdirectory="true"
                  directory="true"
                  multiple
                  onChange={handleFileUpload1}
                  style={{ display: 'none' }}
                />
                <label htmlFor="folderInput1" className="action-btn" title="Add folder">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                  </svg>
                </label>
              </div>
            </div>
            
            {/* Attachments List for First Description */}
            {formData.attachments1.length > 0 && (
              <div className="attachments-list">
                {formData.attachments1.map((file, index) => (
                  <div key={index} className="attachment-item">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                      <polyline points="13 2 13 9 20 9"></polyline>
                    </svg>
                    <span>{file.name}</span>
                    <button type="button" onClick={() => removeAttachment1(index)} className="remove-btn">
                      ×
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Due Date and Category Row */}
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="dueDate">
                Due Date <span className="required">*</span>
              </label>
              <input
                type="date"
                id="dueDate"
                name="dueDate"
                value={formData.dueDate}
                onChange={handleChange}
                className={errors.dueDate ? 'error' : ''}
              />
              {errors.dueDate && <span className="error-message">{errors.dueDate}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="category">Category</label>
              <div className="select-wrapper">
                <select
                  id="category"
                  name="category"
                  value={formData.category}
                  onChange={handleChange}
                >
                  <option value="academic">Academic</option>
                  <option value="personal">Personal</option>
                </select>
                <svg className="select-icon" width="12" height="12" viewBox="0 0 12 12">
                  <path fill="currentColor" d="M6 9L1 4h10z"/>
                </svg>
              </div>
            </div>
          </div>

          {/* Priority Radio Button (High only) */}
          <div className="form-group">
            <label>Priority</label>
            <p className="priority-description">Please mark it as high priority if your task deadline is within 2 days</p>
            <div className="priority-options">
              <label className="priority-label">
                <input 
                  type="checkbox" 
                  name="priority" 
                  value="high"
                  checked={formData.priority === 'high'}
                  onChange={handlePriorityChange}
                />
                <span className="priority-custom">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="10"></circle>
                  </svg>
                  <span className="priority-text">High</span>
                </span>
              </label>
            </div>
          </div>

          {/* Second Description - Priority Description */}
          <div className="form-group">
            <label htmlFor="description2">Priority Description</label>
            <div className="description-wrapper">
              <textarea
                id="description2"
                name="description2"
                value={formData.description2}
                onChange={handleChange}
                placeholder="Add priority-specific details about this task..."
                rows="4"
              />
              <div className="description-actions">
                <input
                  type="file"
                  id="imageInput2"
                  accept="image/*"
                  multiple
                  onChange={handleFileUpload2}
                  style={{ display: 'none' }}
                />
                <label htmlFor="imageInput2" className="action-btn" title="Add image">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                    <circle cx="8.5" cy="8.5" r="1.5"></circle>
                    <polyline points="21 15 16 10 5 21"></polyline>
                  </svg>
                </label>

                <input
                  type="file"
                  id="fileInput2"
                  multiple
                  onChange={handleFileUpload2}
                  style={{ display: 'none' }}
                />
                <label htmlFor="fileInput2" className="action-btn" title="Add file">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                    <polyline points="13 2 13 9 20 9"></polyline>
                  </svg>
                </label>

                <input
                  type="file"
                  id="folderInput2"
                  webkitdirectory="true"
                  directory="true"
                  multiple
                  onChange={handleFileUpload2}
                  style={{ display: 'none' }}
                />
                <label htmlFor="folderInput2" className="action-btn" title="Add folder">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                  </svg>
                </label>
              </div>
            </div>
            
            {/* Attachments List for Second Description */}
            {formData.attachments2.length > 0 && (
              <div className="attachments-list">
                {formData.attachments2.map((file, index) => (
                  <div key={index} className="attachment-item">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                      <polyline points="13 2 13 9 20 9"></polyline>
                    </svg>
                    <span>{file.name}</span>
                    <button type="button" onClick={() => removeAttachment2(index)} className="remove-btn">
                      ×
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Error Message */}
          {errors.submit && (
            <div className="error-banner">
              {errors.submit}
            </div>
          )}

          {/* Success Message */}
          {submitSuccess && (
            <div className="success-banner">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              Task created successfully!
            </div>
          )}

          {/* Action Buttons */}
          <div className="form-actions">
            <button
              type="button"
              onClick={handleReset}
              className="btn-secondary"
              disabled={isSubmitting}
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="1 4 1 10 7 10"></polyline>
                <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
              </svg>
              Reset
            </button>
            <button
              type="submit"
              className="btn-primary"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <svg className="spinner" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <line x1="12" y1="2" x2="12" y2="6"></line>
                    <line x1="12" y1="18" x2="12" y2="22"></line>
                    <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line>
                    <line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line>
                    <line x1="2" y1="12" x2="6" y2="12"></line>
                    <line x1="18" y1="12" x2="22" y2="12"></line>
                    <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line>
                    <line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line>
                  </svg>
                  Creating...
                </>
              ) : (
                <>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                  Create Task
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateTask;
