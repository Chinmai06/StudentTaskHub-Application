import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
<<<<<<< HEAD
import { useAuth } from './context/AuthContext';
=======
import { useAuth } from '../context/AuthContext';
>>>>>>> 51dc39fbfe5540789030f329bde9653cc121e72f
import { useTasks } from '../context/TaskContext';
import { isNonEmptyText } from '../utils/validation';

const initialForm = {
  title: '',
  description: '',
  deadline: '',
  priority: false
};

function CreateTaskPage() {
  const [formData, setFormData] = useState(initialForm);
  const [message, setMessage] = useState('');
  const [errors, setErrors] = useState({});
  const { addTask } = useTasks();
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleChange = (event) => {
    const { name, value, type, checked } = event.target;
    setFormData((current) => ({
      ...current,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setErrors({});
    setMessage('');

    const nextErrors = {};
    if (!isNonEmptyText(formData.title)) nextErrors.title = 'Title is required.';
    if (!isNonEmptyText(formData.description)) nextErrors.description = 'Description is required.';
    if (!formData.deadline) nextErrors.deadline = 'Due date is required.';

    setErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) return;

    try {
      await addTask({
        title: formData.title.trim(),
        description: formData.description.trim(),
        deadline: formData.deadline,
        priority: formData.priority ? 'high' : 'normal',
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

        <div className="two-column-grid">
          <label>
            Due Date
            <input type="date" name="deadline" value={formData.deadline} onChange={handleChange} />
          </label>

          <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '1.5rem' }}>
            <input
              type="checkbox"
              name="priority"
              checked={formData.priority}
              onChange={handleChange}
            />
            High Priority
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
