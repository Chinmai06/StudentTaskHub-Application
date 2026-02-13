# StudentTaskHub - React Version

React application with **TWO description boxes** and **priority selection**!

## ✨ Features

### 📝 Two Description Boxes:
1. **Description** - Appears right after Task Title
2. **Priority Description** - Appears after Priority selection

Each description box has:
- 🖼️ **Image Button** - Upload images
- 📄 **File Button** - Upload any file
- 📁 **Folder Button** - Upload entire folders
- **NO microphone button** (removed as requested)

### 🎯 Form Structure:
1. **Task Title** (required)
2. **Description** (with upload buttons)
3. **Due Date** (required) | **Category** (dropdown)
4. **Priority** - Single "High" checkbox (wider button)
5. **Priority Description** (with upload buttons)
6. **Reset** | **Create Task** buttons

### Other Features:
- ✅ Form validation
- ✅ Success/error messages
- ✅ Separate file lists for each description
- ✅ Mobile responsive
- ✅ Beautiful purple gradient
- ✅ Status field removed
- ✅ Wider priority button

## 🚀 Quick Setup

```bash
# 1. Create React app
npx create-react-app studenttaskhub
cd studenttaskhub

# 2. Copy files to src/ directory:
# - CreateTask.jsx → src/
# - CreateTask.css → src/
# - App.js → src/
# - index.js → src/

# 3. Install dependencies (if needed)
npm install

# 4. Run the app
npm start
```

The app will open at `http://localhost:3000`

## 📁 File Structure

```
studenttaskhub/
├── public/
│   └── index.html
├── src/
│   ├── App.js              # Main app component
│   ├── CreateTask.jsx      # Task creation form
│   ├── CreateTask.css      # All styles
│   └── index.js            # Entry point
├── package.json
└── README.md
```

## 🎯 Component Details

### CreateTask Component

**State Management:**
```javascript
{
  title: '',
  description1: '',        // First description
  description2: '',        // Priority description
  dueDate: '',
  priority: '',            // 'high' or ''
  category: 'academic',
  attachments1: [],        // Files for first description
  attachments2: []         // Files for priority description
}
```

**Key Features:**
- Two independent description boxes
- Two separate attachment lists
- Priority as checkbox (can be checked/unchecked)
- No status field
- Category moved to replace old priority position

## 🎨 Styling

### Priority Button (Wider)
```css
.priority-custom {
  min-width: 150px;
  padding: 0.75rem 2.5rem;
  justify-content: center;
}
```

### Description Boxes
Both description boxes have the same styling:
- Bottom-right positioned buttons
- 3 upload buttons (Image, File, Folder)
- Independent file attachment lists

## 📦 File Handling

### First Description Files
```javascript
attachments1: [File, File, ...]
```

### Priority Description Files
```javascript
attachments2: [File, File, ...]
```

### Form Submission
```javascript
{
  "title": "Assignment",
  "description1": "Main description",
  "description2": "Priority details",
  "dueDate": "2024-03-15",
  "priority": "high",
  "category": "academic",
  "attachments1": ["file1.pdf", "file2.jpg"],
  "attachments2": ["urgent.docx"]
}
```

## 🔧 Customization

### Change Priority Button Width
In `CreateTask.css`:
```css
.priority-custom {
  min-width: 150px;  /* Adjust this value */
}
```

### Add More Priority Options
In `CreateTask.jsx`, add more checkboxes/radio buttons:
```jsx
<label className="priority-label">
  <input type="checkbox" name="priority" value="medium" />
  <span className="priority-custom">Medium</span>
</label>
```

### Modify Upload Button Position
In `CreateTask.css`:
```css
.description-actions {
  bottom: 0.75rem;
  right: 1rem;
}
```

## 🌐 Browser Support

| Feature | Chrome | Safari | Firefox | Edge |
|---------|--------|--------|---------|------|
| Image Upload | ✅ | ✅ | ✅ | ✅ |
| File Upload | ✅ | ✅ | ✅ | ✅ |
| Folder Upload | ✅ | ❌ | ❌ | ✅ |

## 📱 Mobile Responsive

- Description buttons stack properly
- Form fields go full width
- Touch-friendly targets
- Smaller icons on mobile
- Priority button adjusts width

## 🔌 Backend Integration

Use FormData to send files:

```javascript
const formData = new FormData();
formData.append('title', data.title);
formData.append('description1', data.description1);
formData.append('description2', data.description2);
formData.append('dueDate', data.dueDate);
formData.append('priority', data.priority);
formData.append('category', data.category);

// Append files from first description
data.attachments1.forEach(file => {
  formData.append('description1_files', file);
});

// Append files from priority description
data.attachments2.forEach(file => {
  formData.append('description2_files', file);
});

// Send to backend
await fetch('/api/tasks', {
  method: 'POST',
  body: formData
});
```

## 🎉 What's New in v2.0

- ✅ Two separate description boxes
- ✅ Status field removed
- ✅ Category moved to priority position
- ✅ Priority changed to single "High" option
- ✅ Wider priority button (150px min-width)
- ✅ No microphone button
- ✅ Independent file attachments for each description
- ✅ Always visible description boxes

---

**Perfect for students managing tasks with detailed priorities!** ✨
