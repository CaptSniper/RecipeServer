import React, { useState, useEffect } from 'react';
import { fetchRecipes, deleteRecipe } from '../api/recipes';
import { Link } from 'react-router-dom';
import './RecipeList.css';

export default function RecipeList() {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadRecipes = async () => {
      try {
        const data = await fetchRecipes();
        // Normalize the data to always be an array
        if (Array.isArray(data)) {
          setRecipes(data);
        } else if (Array.isArray(data.recipes)) {
          setRecipes(data.recipes);
        } else {
          setRecipes([]);
        }
      } catch (err) {
        console.error('Failed to fetch recipes:', err);
        setRecipes([]);
      } finally {
        setLoading(false);
      }
    };
    loadRecipes();
  }, []);

  const handleDelete = async (id) => {
    try {
      await deleteRecipe(id);
      setRecipes(prev => prev.filter(r => r.id !== id));
    } catch (err) {
      console.error('Failed to delete recipe:', err);
      alert('Failed to delete recipe. Check console for details.');
    }
  };

  if (loading) return <p>Loading recipes...</p>;

  if (!recipes.length) return (
    <div className="recipe-list-container">
      <h1>Digital Cookbook</h1>
      <p>No recipes found.</p>
    </div>
  );

  return (
    <div className="recipe-list-container">
      <h1>Digital Cookbook</h1>
      <ul className="recipe-list">
        {recipes.map(r => (
          <li key={r.id} className="recipe-item">
            <Link to={`/recipe/${r.id}`} className="recipe-link">
              <span className="recipe-name">{r.name}</span>
            </Link>
            <button
              onClick={() => handleDelete(r.id)}
              className="delete-button"
              aria-label={`Delete ${r.name}`}
            >
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
