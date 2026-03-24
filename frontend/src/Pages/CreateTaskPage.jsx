import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import { useTasks } from './context/TaskContext';
import { isNonEmptyText } from './utils/validation';

const initialForm = {
  title: '',
  category: 'Study',
  description: '',
  location: '',
  priority: 'Medium',
  dueDate: ''
};

function CreateTaskPage() {
  const [formData, setFormData] = useState(initialForm);
  const [message, setMessage] = useState('');
  const [errors, setErrors] = useState({});
  const { addTask } = useTasks();
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData((currentData) => ({
      ...currentData,
      [name]: value
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const nextErrors = {};

    if (!isNonEmptyText(formData.title)) nextErrors.title = 'Title is required.';
    if (!isNonEmptyText(formData.description)) nextErrors.description = 'Description is required.';
    if (!isNonEmptyText(formData.location)) nextErrors.location = 'Location is required.';
    if (!formData.dueDate) nextErrors.dueDate = 'Due date is required.';

    setErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    addTask(formData, user);
    setFormData(initialForm);
    setMessage('Task created successfully.');
    navigate('/home');
  };

  return (
    <section className="content-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Create</p>
          <h2>Create a New Task</h2>
        </div>
      </div>

      <form className="form-grid" onSubmit={handleSubmit} data-testid="create-task-form">
        <label>
          Task Title
          <input name="title" value={formData.title} onChange={handleChange} placeholder="Enter task title" />
        </label>
        {errors.title ? <p className="error-text">{errors.title}</p> : null}

        <label>
          Category
          <select name="category" value={formData.category} onChange={handleChange}>
            <option>Study</option>
            <option>Project</option>
            <option>Errand</option>
            <option>Event</option>
          </select>
        </label>

        <label>
          Description
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows="4"
            placeholder="Describe the task"
          />
        </label>
        {errors.description ? <p className="error-text">{errors.description}</p> : null}

        <label>
          Location
          <input name="location" value={formData.location} onChange={handleChange} placeholder="Example: Reitz Union" />
        </label>
        {errors.location ? <p className="error-text">{errors.location}</p> : null}

        <div className="two-column-grid">
          <label>
            Priority
            <select name="priority" value={formData.priority} onChange={handleChange}>
              <option>High</option>
              <option>Medium</option>
              <option>Low</option>
            </select>
          </label>

          <label>
            Due Date
            <input type="date" name="dueDate" value={formData.dueDate} onChange={handleChange} />
          </label>
        </div>
        {errors.dueDate ? <p className="error-text">{errors.dueDate}</p> : null}

        <button type="submit" className="primary-button">
          Create Task
        </button>
        {message ? <p className="success-text">{message}</p> : null}
      </form>
    </section>
  );
}

export default CreateTaskPage;
