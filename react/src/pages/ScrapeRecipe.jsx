import React, { useState, useEffect } from 'react';
import { scrapeRecipe, createRecipe, updateRecipe } from '../api/recipes';
import { useNavigate, useParams } from 'react-router-dom';
import './ScrapeRecipe.css';

export default function ScrapeRecipe() {
  const { id } = useParams(); // optional if editing existing recipe
  const navigate = useNavigate();

  const [url, setUrl] = useState('');
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(false);
  const [customProps, setCustomProps] = useState([{ key: '', value: '' }]);

  // Auto-resize textareas
  useEffect(() => {
    const textareas = document.querySelectorAll('.auto-resize');
    textareas.forEach(textarea => {
      textarea.style.height = 'auto';
      textarea.style.height = textarea.scrollHeight + 'px';
    });
  }, [recipe]);

  const handleScrape = async () => {
    if (!url) return;
    setLoading(true);
    try {
      const scraped = await scrapeRecipe(url, false);
      setRecipe(scraped);

      // Initialize custom props from CoreProps if they exist
      if (scraped.CoreProps && Object.keys(scraped.CoreProps).length > 0) {
        const props = Object.entries(scraped.CoreProps).map(([key, value]) => ({ key, value }));
        setCustomProps(props);
      }
    } catch (err) {
      alert('Failed to scrape recipe');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!recipe) return;

    // Build CoreProps from custom properties
    const coreProps = {};
    customProps.forEach(prop => {
      if (prop.key && prop.value) {
        coreProps[prop.key] = prop.value;
      }
    });

    const recipeToSave = {
      ...recipe,
      CoreProps: coreProps
    };

    try {
      if (id) {
        await updateRecipe(id, recipeToSave);
      } else {
        await createRecipe(recipeToSave);
      }
      navigate('/');
    } catch (err) {
      alert('Failed to save recipe');
      console.error(err);
    }
  };

  const updateField = (field, value, index = null) => {
    if (!recipe) return;
    const copy = { ...recipe };
    if (field === 'Name' || field === 'ImagePath') {
      copy[field] = value;
    } else if (field === 'Ingredients') {
      const ing = [...copy.Ingredients];
      ing[index] = value;
      copy.Ingredients = ing;
    } else if (field === 'Steps') {
      const st = [...copy.Steps];
      st[index] = value;
      copy.Steps = st;
    }
    setRecipe(copy);
  };

  const addField = (field) => {
    if (!recipe) return;
    const copy = { ...recipe };
    if (field === 'Ingredients') copy.Ingredients.push('');
    if (field === 'Steps') copy.Steps.push('');
    setRecipe(copy);
  };

  const removeIngredient = (index) => {
    if (!recipe || recipe.Ingredients.length <= 1) return;
    const copy = { ...recipe };
    copy.Ingredients = copy.Ingredients.filter((_, i) => i !== index);
    setRecipe(copy);
  };

  const removeStep = (index) => {
    if (!recipe || recipe.Steps.length <= 1) return;
    const copy = { ...recipe };
    copy.Steps = copy.Steps.filter((_, i) => i !== index);
    setRecipe(copy);
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

  const handleTextareaResize = (e) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';
  };

  return (
    <div className="recipe-form-container">
      {!recipe ? (
        <div className="scrape-form">
          <h1>{id ? 'Edit Recipe' : 'Scrape Recipe'}</h1>
          <p className="scrape-description">Enter a recipe URL to automatically extract recipe information</p>

          <div className="scrape-input-section">
            <div className="info-row">
              <label className="info-label">Recipe URL</label>
              <input
                type="text"
                className="info-value"
                placeholder="https://www.allrecipes.com/recipe/..."
                value={url}
                onChange={e => setUrl(e.target.value)}
              />
            </div>

            <button
              onClick={handleScrape}
              disabled={loading}
              className="primary submit-button"
            >
              {loading ? 'Scraping...' : 'Scrape Recipe'}
            </button>
          </div>
        </div>
      ) : (
        <form onSubmit={e => { e.preventDefault(); handleSave(); }} className="recipe-form">
          <h1>{id ? 'Edit Recipe' : 'Scraped Recipe'}</h1>

          {/* Recipe Information Section */}
          <section className="recipe-info-section">
            <h3>Recipe Information</h3>

            <div className="info-row">
              <label className="info-label">Name</label>
              <input
                className="info-value"
                placeholder="Enter recipe name"
                value={recipe.Name}
                onChange={e => updateField('Name', e.target.value)}
                required
              />
            </div>

            <div className="info-row">
              <label className="info-label">Image URL</label>
              <input
                className="info-value"
                placeholder="Enter image URL (optional)"
                value={recipe.ImagePath || ''}
                onChange={e => updateField('ImagePath', e.target.value)}
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
              {recipe.Ingredients.map((ing, i) => (
                <div key={i} className="info-row">
                  <label className="info-label">{i + 1}</label>
                  <div className="custom-prop-value-container">
                    <textarea
                      className="info-value auto-resize"
                      placeholder="Enter ingredient"
                      value={ing}
                      rows={1}
                      onInput={handleTextareaResize}
                      onChange={e => updateField('Ingredients', e.target.value, i)}
                    />
                    {recipe.Ingredients.length > 1 && (
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
            <button type="button" className="add-button" onClick={() => addField('Ingredients')}>
              + Add Ingredient
            </button>
          </section>

          {/* Steps Section */}
          <section className="recipe-section">
            <h3>Steps</h3>
            <div className="items-list">
              {recipe.Steps.map((step, i) => (
                <div key={i} className="info-row">
                  <label className="info-label">{i + 1}</label>
                  <div className="custom-prop-value-container">
                    <textarea
                      className="info-value auto-resize"
                      placeholder="Enter step"
                      value={step}
                      rows={1}
                      onInput={handleTextareaResize}
                      onChange={e => updateField('Steps', e.target.value, i)}
                    />
                    {recipe.Steps.length > 1 && (
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
            <button type="button" className="add-button" onClick={() => addField('Steps')}>
              + Add Step
            </button>
          </section>

          <button type="submit" className="primary submit-button">
            {id ? 'Update Recipe' : 'Save Recipe'}
          </button>
        </form>
      )}
    </div>
  );
}
