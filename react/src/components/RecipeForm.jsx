import React, { useState, useEffect } from 'react';
import { createRecipe } from '../api/recipes';
import { useNavigate } from 'react-router-dom';
import './RecipeForm.css';

export default function RecipeForm() {
  const [name, setName] = useState('');
  const [ingredients, setIngredients] = useState(['']);
  const [steps, setSteps] = useState(['']);
  const [customProps, setCustomProps] = useState([{ key: '', value: '' }]);
  const navigate = useNavigate();

  // Auto-resize textareas
  useEffect(() => {
    const textareas = document.querySelectorAll('.auto-resize');
    textareas.forEach(textarea => {
      textarea.style.height = 'auto';
      textarea.style.height = textarea.scrollHeight + 'px';
    });
  }, [ingredients, steps]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Build CoreProps from custom properties
    const coreProps = {};
    customProps.forEach(prop => {
      if (prop.key && prop.value) {
        coreProps[prop.key] = prop.value;
      }
    });

    await createRecipe({
      Name: name,
      Ingredients: ingredients.filter(Boolean),
      Steps: steps.filter(Boolean),
      CoreProps: coreProps
    });
    navigate('/');
  };

  const updateCustomProp = (index, field, value) => {
    const copy = [...customProps];
    copy[index][field] = value;
    setCustomProps(copy);
  };

  const addCustomProp = () => {
    setCustomProps([...customProps, { key: '', value: '' }]);
  };

  const removeCustomProp = (index) => {
    setCustomProps(customProps.filter((_, i) => i !== index));
  };

  const removeIngredient = (index) => {
    if (ingredients.length > 1) {
      setIngredients(ingredients.filter((_, i) => i !== index));
    }
  };

  const removeStep = (index) => {
    if (steps.length > 1) {
      setSteps(steps.filter((_, i) => i !== index));
    }
  };

  const handleTextareaResize = (e) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';
  };

  return (
    <div className="recipe-form-container">
      <form onSubmit={handleSubmit} className="recipe-form">
        <h1>New Recipe</h1>

        {/* Recipe Information Section */}
        <section className="recipe-info-section">
          <h3>Recipe Information</h3>

          <div className="info-row">
            <label className="info-label">Name</label>
            <input
              className="info-value"
              placeholder="Enter recipe name"
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
          </div>

          {/* Custom Properties */}
          {customProps.map((prop, i) => (
            <div key={i} className="info-row custom-prop-row">
              <input
                className="info-label-input"
                placeholder="e.g. Cook Time"
                value={prop.key}
                onChange={e => updateCustomProp(i, 'key', e.target.value)}
              />
              <div className="custom-prop-value-container">
                <input
                  className="info-value"
                  placeholder="e.g. 30 minutes"
                  value={prop.value}
                  onChange={e => updateCustomProp(i, 'value', e.target.value)}
                />
                {customProps.length > 1 && (
                  <button
                    type="button"
                    className="remove-button"
                    onClick={() => removeCustomProp(i)}
                    aria-label="Remove property"
                  >
                    ×
                  </button>
                )}
              </div>
            </div>
          ))}

          <button type="button" className="add-button" onClick={addCustomProp}>
            + Add Property
          </button>
        </section>

        {/* Ingredients Section */}
        <section className="recipe-section">
          <h3>Ingredients</h3>
          <div className="items-list">
            {ingredients.map((ing, i) => (
              <div key={i} className="info-row">
                <label className="info-label">{i + 1}</label>
                <div className="custom-prop-value-container">
                  <textarea
                    className="info-value auto-resize"
                    placeholder="Enter ingredient"
                    value={ing}
                    rows={1}
                    onInput={handleTextareaResize}
                    onChange={e => {
                      const copy = [...ingredients];
                      copy[i] = e.target.value;
                      setIngredients(copy);
                    }}
                  />
                  {ingredients.length > 1 && (
                    <button
                      type="button"
                      className="remove-button"
                      onClick={() => removeIngredient(i)}
                      aria-label="Remove ingredient"
                    >
                      ×
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
          <button type="button" className="add-button" onClick={() => setIngredients([...ingredients, ''])}>
            + Add Ingredient
          </button>
        </section>

        {/* Steps Section */}
        <section className="recipe-section">
          <h3>Steps</h3>
          <div className="items-list">
            {steps.map((step, i) => (
              <div key={i} className="info-row">
                <label className="info-label">{i + 1}</label>
                <div className="custom-prop-value-container">
                  <textarea
                    className="info-value auto-resize"
                    placeholder="Enter step"
                    value={step}
                    rows={1}
                    onInput={handleTextareaResize}
                    onChange={e => {
                      const copy = [...steps];
                      copy[i] = e.target.value;
                      setSteps(copy);
                    }}
                  />
                  {steps.length > 1 && (
                    <button
                      type="button"
                      className="remove-button"
                      onClick={() => removeStep(i)}
                      aria-label="Remove step"
                    >
                      ×
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
          <button type="button" className="add-button" onClick={() => setSteps([...steps, ''])}>
            + Add Step
          </button>
        </section>

        <button type="submit" className="primary submit-button">Create Recipe</button>
      </form>
    </div>
  );
}
