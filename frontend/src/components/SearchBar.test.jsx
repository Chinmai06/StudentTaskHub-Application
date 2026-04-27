//frontend/src/components/SearchBar.test.jsx
import { render, screen, fireEvent } from '@testing-library/react';
import SearchBar from './SearchBar';

const defaultProps = {
  searchTerm: '',
  setSearchTerm: vi.fn(),
  statusFilter: '',
  setStatusFilter: vi.fn(),
  priorityFilter: '',
  setPriorityFilter: vi.fn(),
  sortBy: '',
  setSortBy: vi.fn(),
  onFilter: vi.fn(),
  onClear: vi.fn(),
};

describe('SearchBar', () => {
  it('renders the search input', () => {
    render(<SearchBar {...defaultProps} />);
    expect(screen.getByPlaceholderText('Search tasks...')).toBeInTheDocument();
  });

  it('calls setSearchTerm when typing', () => {
    render(<SearchBar {...defaultProps} />);
    fireEvent.change(screen.getByPlaceholderText('Search tasks...'), {
      target: { value: 'math' },
    });
    expect(defaultProps.setSearchTerm).toHaveBeenCalledWith('math');
  });

  it('calls onFilter when Enter is pressed', () => {
    render(<SearchBar {...defaultProps} />);
    fireEvent.keyDown(screen.getByPlaceholderText('Search tasks...'), { key: 'Enter' });
    expect(defaultProps.onFilter).toHaveBeenCalled();
  });

  it('calls onFilter when Filter button is clicked', () => {
    render(<SearchBar {...defaultProps} />);
    fireEvent.click(screen.getByText('Filter'));
    expect(defaultProps.onFilter).toHaveBeenCalled();
  });

  it('calls onClear when Clear button is clicked', () => {
    render(<SearchBar {...defaultProps} />);
    fireEvent.click(screen.getByText('Clear'));
    expect(defaultProps.onClear).toHaveBeenCalled();
  });

  it('calls setStatusFilter on status change', () => {
    render(<SearchBar {...defaultProps} />);
    fireEvent.change(screen.getByLabelText('Filter by status'), {
      target: { value: 'open' },
    });
    expect(defaultProps.setStatusFilter).toHaveBeenCalledWith('open');
  });
});
