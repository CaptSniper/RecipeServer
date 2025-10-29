import React, { useState } from 'react';
import { scrapeRecipe, createRecipe, updateRecipe } from '../api/recipes';
import { useNavigate, useParams } from 'react-router-dom';

export default function ScrapeRecipe() {
  const { id } = useParams(); // optional if editing existing recipe
  const navigate = useNavigate();

  const [url, setUrl] = useState('');
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleScrape = async () => {
    if (!url) return;
    setLoading(true);
    try {
      const scraped = await scrapeRecipe(url, false);
      setRecipe(scraped);
    } catch (err) {
      alert('Failed to scrape recipe');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!recipe) return;
    try {
      if (id) {
        await updateRecipe(id, recipe);
      } else {
        await createRecipe(recipe);
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

  return (
    <div>
      <button onClick={() => navigate('/recipes')}>Back to cookbook</button>
      <h1>{id ? 'Edit Recipe' : 'Scrape Recipe'}</h1>
      {!recipe && (
        <div>
          <input
            type="text"
            placeholder="Enter recipe URL"
            value={url}
            onChange={e => setUrl(e.target.value)}
          />
          <button onClick={handleScrape} disabled={loading}>
            {loading ? 'Scraping...' : 'Scrape'}
          </button>
        </div>
      )}

      {recipe && (
        <form onSubmit={e => { e.preventDefault(); handleSave(); }}>
          <input
            placeholder="Recipe Name"
            value={recipe.Name}
            onChange={e => updateField('Name', e.target.value)}
            required
          />
          <input
            placeholder="Image URL"
            value={recipe.ImagePath || ''}
            onChange={e => updateField('ImagePath', e.target.value)}
          />

          <h2>Ingredients</h2>
          {recipe.Ingredients.map((ing, i) => (
            <input
              key={i}
              value={ing}
              onChange={e => updateField('Ingredients', e.target.value, i)}
            />
          ))}
          <button type="button" onClick={() => addField('Ingredients')}>Add Ingredient</button>

          <h2>Steps</h2>
          {recipe.Steps.map((step, i) => (
            <input
              key={i}
              value={step}
              onChange={e => updateField('Steps', e.target.value, i)}
            />
          ))}
          <button type="button" onClick={() => addField('Steps')}>Add Step</button>

          <br />
          <button type="submit">{id ? 'Update Recipe' : 'Save Recipe'}</button>
        </form>
      )}
    </div>
  );
}
