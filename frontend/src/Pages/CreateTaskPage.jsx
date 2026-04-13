import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTasks } from '../context/TaskContext';
import { isNonEmptyText } from '../utils/validation';

const initialForm = {
  title: '',
  category: 'Study',
  description: '',
  location: '',
  priority: 'Medium',
  deadline: ''
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
    setFormData((current) => ({
      ...current,
      [name]: value
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setErrors({});
    setMessage('');

    const nextErrors = {};
    if (!isNonEmptyText(formData.title)) nextErrors.title = 'Title is required.';
    if (!isNonEmptyText(formData.description)) nextErrors.description = 'Description is required.';
    if (!isNonEmptyText(formData.location)) nextErrors.location = 'Location is required.';
    if (!formData.deadline) nextErrors.deadline = 'Due date is required.';

    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) return;

    try {
      await addTask({
        title: formData.title.trim(),
        category: formData.category,
        description: formData.description.trim(),
        location: formData.location.trim(),
        priority: formData.priority,
        deadline: formData.deadline,
        created_by: user.username
      });
      setFormData(initialForm);
      setMessage('Task created successfully.');
      setTimeout(() => navigate('/home'), 1000);
    } catch (err) {
      setErrors({ form: err.message });
    }
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
            <input type="date" name="deadline" value={formData.deadline} onChange={handleChange} />
          </label>
        </div>
        {errors.deadline ? <p className="error-text">{errors.deadline}</p> : null}

        {errors.form ? <p className="error-text">{errors.form}</p> : null}
        {message ? <p className="success-text">{message}</p> : null}

        <button type="submit" className="primary-button">
          Create Task
        </button>
      </form>
    </section>
  );
}

export default CreateTaskPage;
