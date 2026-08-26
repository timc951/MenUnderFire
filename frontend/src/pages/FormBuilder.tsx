import { useState, useEffect, useCallback, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import { FormDetail, FormField, FormFieldType } from '../types';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { LoadingSpinner } from '../components/common/LoadingSpinner';

const FIELD_TYPE_LABELS: Record<FormFieldType, string> = {
  TEXT_DISPLAY: 'Text Display (Read-only)',
  TEXT_SMALL: 'Small Text (1 line)',
  TEXT_MEDIUM: 'Medium Text (3 lines)',
  TEXT_LARGE: 'Large Text (6 lines)',
  CHECKBOX: 'Checkboxes (Multiple Selection)',
  RADIO: 'Radio Buttons (Single Selection)',
  DROPDOWN: 'Dropdown (Single Selection)',
};

const FIELD_TYPES: FormFieldType[] = [
  'TEXT_DISPLAY',
  'TEXT_SMALL',
  'TEXT_MEDIUM',
  'TEXT_LARGE',
  'CHECKBOX',
  'RADIO',
  'DROPDOWN',
];

export function FormBuilder() {
  const { formId } = useParams<{ formId: string }>();
  const navigate = useNavigate();
  const api = useApi();
  const [form, setForm] = useState<FormDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showAddField, setShowAddField] = useState(false);
  const [editingField, setEditingField] = useState<FormField | null>(null);

  // New field state
  const [fieldType, setFieldType] = useState<FormFieldType>('TEXT_SMALL');
  const [fieldLabel, setFieldLabel] = useState('');
  const [fieldDescription, setFieldDescription] = useState('');
  const [fieldRequired, setFieldRequired] = useState(false);
  const [fieldOptions, setFieldOptions] = useState<string[]>(['']);
  const [isSaving, setIsSaving] = useState(false);
  const [draggedFieldId, setDraggedFieldId] = useState<string | null>(null);
  const [dragOverFieldId, setDragOverFieldId] = useState<string | null>(null);
  const optionInputRefs = useRef<(HTMLInputElement | null)[]>([]);

  const fetchForm = useCallback(async () => {
    if (!formId) return;
    setIsLoading(true);
    try {
      const data = await api.get<FormDetail>(`/forms/${formId}`);
      setForm(data);
    } catch (err) {
      console.error('Failed to load form:', err);
    } finally {
      setIsLoading(false);
    }
  }, [api, formId]);

  useEffect(() => {
    fetchForm();
  }, [fetchForm]);

  const resetFieldForm = () => {
    setFieldType('TEXT_SMALL');
    setFieldLabel('');
    setFieldDescription('');
    setFieldRequired(false);
    setFieldOptions(['']);
    setEditingField(null);
  };

  const handleAddField = async () => {
    if (!fieldLabel.trim() || !formId) return;

    const needsOptions = ['CHECKBOX', 'RADIO', 'DROPDOWN'].includes(fieldType);
    const validOptions = fieldOptions.filter((o) => o.trim());
    if (needsOptions && validOptions.length === 0) {
      alert('Please add at least one option');
      return;
    }

    setIsSaving(true);
    try {
      await api.post(`/forms/${formId}/fields`, {
        fieldType,
        label: fieldLabel.trim(),
        description: fieldDescription.trim() || undefined,
        isRequired: fieldRequired,
        options: needsOptions ? validOptions : undefined,
      });
      setShowAddField(false);
      resetFieldForm();
      fetchForm();
    } catch (err) {
      console.error('Failed to add field:', err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleUpdateField = async () => {
    if (!fieldLabel.trim() || !editingField) return;

    const needsOptions = ['CHECKBOX', 'RADIO', 'DROPDOWN'].includes(editingField.fieldType);
    const validOptions = fieldOptions.filter((o) => o.trim());
    if (needsOptions && validOptions.length === 0) {
      alert('Please add at least one option');
      return;
    }

    setIsSaving(true);
    try {
      await api.put(`/forms/${formId}/fields/${editingField.id}`, {
        label: fieldLabel.trim(),
        description: fieldDescription.trim() || undefined,
        isRequired: fieldRequired,
        options: needsOptions ? validOptions : undefined,
      });
      setEditingField(null);
      resetFieldForm();
      fetchForm();
    } catch (err) {
      console.error('Failed to update field:', err);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteField = async (fieldId: string) => {
    if (!confirm('Are you sure you want to delete this field?')) return;
    try {
      await api.delete(`/forms/${formId}/fields/${fieldId}`);
      fetchForm();
    } catch (err) {
      console.error('Failed to delete field:', err);
    }
  };

  const startEditingField = (field: FormField) => {
    setEditingField(field);
    setFieldLabel(field.label);
    setFieldDescription(field.description || '');
    setFieldRequired(field.isRequired);
    setFieldOptions(field.options?.length ? field.options : ['']);
    setShowAddField(false);
  };

  const addOption = () => {
    setFieldOptions([...fieldOptions, '']);
    // Focus the new input after React renders it
    setTimeout(() => {
      const newIndex = fieldOptions.length;
      optionInputRefs.current[newIndex]?.focus();
    }, 0);
  };

  const handleOptionKeyDown = (e: React.KeyboardEvent<HTMLInputElement>, _index: number) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addOption();
    }
  };

  const updateOption = (index: number, value: string) => {
    const newOptions = [...fieldOptions];
    newOptions[index] = value;
    setFieldOptions(newOptions);
  };

  const removeOption = (index: number) => {
    if (fieldOptions.length <= 1) return;
    setFieldOptions(fieldOptions.filter((_, i) => i !== index));
  };

  // Drag and drop handlers
  const handleDragStart = (fieldId: string) => {
    setDraggedFieldId(fieldId);
  };

  const handleDragOver = (e: React.DragEvent, fieldId: string) => {
    e.preventDefault();
    if (fieldId !== draggedFieldId) {
      setDragOverFieldId(fieldId);
    }
  };

  const handleDragLeave = () => {
    setDragOverFieldId(null);
  };

  const handleDrop = async (targetFieldId: string) => {
    if (!draggedFieldId || !form || draggedFieldId === targetFieldId) {
      setDraggedFieldId(null);
      setDragOverFieldId(null);
      return;
    }

    // Reorder fields locally first for immediate feedback
    const fields = [...form.fields];
    const draggedIndex = fields.findIndex((f) => f.id === draggedFieldId);
    const targetIndex = fields.findIndex((f) => f.id === targetFieldId);

    if (draggedIndex === -1 || targetIndex === -1) return;

    const [draggedField] = fields.splice(draggedIndex, 1);
    fields.splice(targetIndex, 0, draggedField);

    // Update local state immediately
    setForm({ ...form, fields });
    setDraggedFieldId(null);
    setDragOverFieldId(null);

    // Send reorder to backend
    try {
      await api.put(`/forms/${formId}/fields/reorder`, {
        fieldIds: fields.map((f) => f.id),
      });
    } catch (err) {
      console.error('Failed to reorder fields:', err);
      // Refresh to get correct order on error
      fetchForm();
    }
  };

  const handleDragEnd = () => {
    setDraggedFieldId(null);
    setDragOverFieldId(null);
  };

  if (isLoading) {
    return <LoadingSpinner size="lg" />;
  }

  if (!form) {
    return <p className="text-red-600 dark:text-red-400">Form not found</p>;
  }

  const needsOptions = ['CHECKBOX', 'RADIO', 'DROPDOWN'].includes(fieldType);
  const editingNeedsOptions = editingField && ['CHECKBOX', 'RADIO', 'DROPDOWN'].includes(editingField.fieldType);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={() => navigate(`/organizations/${form.organizationId}/forms`)}
            className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            aria-label="Back"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{form.name}</h1>
            {form.description && (
              <p className="text-gray-600 dark:text-stone-400">{form.description}</p>
            )}
          </div>
        </div>
        <Button onClick={() => { setShowAddField(true); resetFieldForm(); }}>
          Add Field
        </Button>
      </div>

      {/* Field List */}
      {form.fields.length === 0 ? (
        <Card>
          <p className="text-gray-500 dark:text-stone-400 text-center py-8">
            No fields yet. Add your first field to start building the form.
          </p>
        </Card>
      ) : (
        <div className="space-y-3">
          {form.fields.map((field, index) => (
            <div
              key={field.id}
              draggable
              onDragStart={() => handleDragStart(field.id)}
              onDragOver={(e) => handleDragOver(e, field.id)}
              onDragLeave={handleDragLeave}
              onDrop={() => handleDrop(field.id)}
              onDragEnd={handleDragEnd}
              className={`transition-all ${
                draggedFieldId === field.id ? 'opacity-50' : ''
              } ${
                dragOverFieldId === field.id ? 'ring-2 ring-amber-500 ring-offset-2 dark:ring-offset-stone-900' : ''
              }`}
            >
              <Card>
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3 cursor-grab active:cursor-grabbing">
                    <svg className="w-5 h-5 text-gray-400 dark:text-stone-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" />
                    </svg>
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-medium px-2 py-0.5 rounded bg-stone-100 dark:bg-stone-700 text-stone-600 dark:text-stone-400">
                          {index + 1}
                        </span>
                    <span className="text-xs font-medium px-2 py-0.5 rounded bg-amber-100 dark:bg-amber-900/50 text-amber-700 dark:text-amber-400">
                      {FIELD_TYPE_LABELS[field.fieldType]}
                    </span>
                    {field.isRequired && (
                      <span className="text-xs font-medium px-2 py-0.5 rounded bg-red-100 dark:bg-red-900/50 text-red-700 dark:text-red-400">
                        Required
                      </span>
                    )}
                  </div>
                  <h3 className="font-medium text-gray-900 dark:text-white mt-2">{field.label}</h3>
                  {field.description && (
                    <p className="text-sm text-gray-600 dark:text-stone-400 mt-1">{field.description}</p>
                  )}
                  {field.options && field.options.length > 0 && (
                    <div className="mt-2 flex flex-wrap gap-2">
                      {field.options.map((option, i) => (
                        <span
                          key={i}
                          className="text-xs px-2 py-1 bg-stone-100 dark:bg-stone-700 text-stone-700 dark:text-stone-300 rounded"
                        >
                          {option}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
                  </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => startEditingField(field)}
                    className="p-2 text-gray-400 hover:text-amber-500"
                    title="Edit field"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    onClick={() => handleDeleteField(field.id)}
                    className="p-2 text-gray-400 hover:text-red-500"
                    title="Delete field"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
            </Card>
            </div>
          ))}
        </div>
      )}

      {/* Add Field Modal */}
      {showAddField && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-lg mx-4 my-8">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Add New Field</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Field Type
                </label>
                <select
                  value={fieldType}
                  onChange={(e) => setFieldType(e.target.value as FormFieldType)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                >
                  {FIELD_TYPES.map((type) => (
                    <option key={type} value={type}>
                      {FIELD_TYPE_LABELS[type]}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Label
                </label>
                <input
                  type="text"
                  value={fieldLabel}
                  onChange={(e) => setFieldLabel(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter field label"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Description (optional)
                </label>
                <input
                  type="text"
                  value={fieldDescription}
                  onChange={(e) => setFieldDescription(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter description"
                />
              </div>

              {fieldType !== 'TEXT_DISPLAY' && (
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="field-required"
                    checked={fieldRequired}
                    onChange={(e) => setFieldRequired(e.target.checked)}
                    className="h-4 w-4 rounded border-gray-300 dark:border-stone-600 text-amber-600 focus:ring-amber-500"
                  />
                  <label htmlFor="field-required" className="text-sm text-gray-700 dark:text-stone-300">
                    Required field
                  </label>
                </div>
              )}

              {needsOptions && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-2">
                    Options
                  </label>
                  <div className="space-y-2">
                    {fieldOptions.map((option, index) => (
                      <div key={index} className="flex gap-2">
                        <input
                          ref={(el) => { optionInputRefs.current[index] = el; }}
                          type="text"
                          value={option}
                          onChange={(e) => updateOption(index, e.target.value)}
                          onKeyDown={(e) => handleOptionKeyDown(e, index)}
                          className="flex-1 px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                          placeholder={`Option ${index + 1}`}
                        />
                        <button
                          type="button"
                          onClick={() => removeOption(index)}
                          className="p-2 text-gray-400 hover:text-red-500"
                          disabled={fieldOptions.length <= 1}
                        >
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    ))}
                  </div>
                  <Button variant="secondary" onClick={addOption} className="mt-2">
                    Add Option
                  </Button>
                </div>
              )}
            </div>

            <div className="flex justify-end gap-3 mt-6">
              <Button variant="secondary" onClick={() => { setShowAddField(false); resetFieldForm(); }}>
                Cancel
              </Button>
              <Button onClick={handleAddField} disabled={isSaving || !fieldLabel.trim()}>
                {isSaving ? 'Adding...' : 'Add Field'}
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Field Modal */}
      {editingField && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 overflow-y-auto">
          <div className="bg-white dark:bg-stone-800 rounded-lg p-6 w-full max-w-lg mx-4 my-8">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Edit Field</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Field Type
                </label>
                <p className="px-3 py-2 bg-stone-100 dark:bg-stone-700 text-gray-700 dark:text-stone-300 rounded-md">
                  {FIELD_TYPE_LABELS[editingField.fieldType]}
                </p>
                <p className="text-xs text-gray-500 dark:text-stone-500 mt-1">
                  Field type cannot be changed after creation
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Label
                </label>
                <input
                  type="text"
                  value={fieldLabel}
                  onChange={(e) => setFieldLabel(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter field label"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
                  Description (optional)
                </label>
                <input
                  type="text"
                  value={fieldDescription}
                  onChange={(e) => setFieldDescription(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="Enter description"
                />
              </div>

              {editingField.fieldType !== 'TEXT_DISPLAY' && (
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="edit-field-required"
                    checked={fieldRequired}
                    onChange={(e) => setFieldRequired(e.target.checked)}
                    className="h-4 w-4 rounded border-gray-300 dark:border-stone-600 text-amber-600 focus:ring-amber-500"
                  />
                  <label htmlFor="edit-field-required" className="text-sm text-gray-700 dark:text-stone-300">
                    Required field
                  </label>
                </div>
              )}

              {editingNeedsOptions && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-2">
                    Options
                  </label>
                  <div className="space-y-2">
                    {fieldOptions.map((option, index) => (
                      <div key={index} className="flex gap-2">
                        <input
                          ref={(el) => { optionInputRefs.current[index] = el; }}
                          type="text"
                          value={option}
                          onChange={(e) => updateOption(index, e.target.value)}
                          onKeyDown={(e) => handleOptionKeyDown(e, index)}
                          className="flex-1 px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
                          placeholder={`Option ${index + 1}`}
                        />
                        <button
                          type="button"
                          onClick={() => removeOption(index)}
                          className="p-2 text-gray-400 hover:text-red-500"
                          disabled={fieldOptions.length <= 1}
                        >
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    ))}
                  </div>
                  <Button variant="secondary" onClick={addOption} className="mt-2">
                    Add Option
                  </Button>
                </div>
              )}
            </div>

            <div className="flex justify-end gap-3 mt-6">
              <Button variant="secondary" onClick={() => { setEditingField(null); resetFieldForm(); }}>
                Cancel
              </Button>
              <Button onClick={handleUpdateField} disabled={isSaving || !fieldLabel.trim()}>
                {isSaving ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
