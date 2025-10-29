import React, { useState, useEffect } from 'react';
import { fetchRecipes, deleteRecipe } from '../api/recipes';
import { useNavigate } from 'react-router-dom';

export default function RecipeList() {
  const [recipes, setRecipes] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

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
    <div>
      <h1>Recipes</h1>
      <p>No recipes found.</p>
      <button onClick={() => navigate('/new')}>New Recipe</button>
    </div>
  );

  return (
    <div>
      <h1>Recipes</h1>
      <button onClick={() => navigate('/new')}>New Recipe</button>
      <ul>
        {recipes.map(r => (
          <li key={r.id}>
            <span
              onClick={() => navigate(`/recipe/${r.id}`)}
              style={{ cursor: 'pointer', textDecoration: 'underline' }}
            >
              {r.name}
            </span>
            <button onClick={() => handleDelete(r.id)} style={{ marginLeft: '10px' }}>
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
