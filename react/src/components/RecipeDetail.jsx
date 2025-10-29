import React, { useEffect, useState } from 'react';
import { fetchRecipe } from '../api/recipes';
import { useParams, useNavigate } from 'react-router-dom';

export default function RecipeDetail() {
  const { id } = useParams();
  const [recipe, setRecipe] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchRecipe(id).then(setRecipe);
  }, [id]);

  if (!recipe) return <div>Loading...</div>;

  return (
    <div>
      <button onClick={() => navigate(-1)}>Back</button>
      <h1>{recipe.Name}</h1>
      {recipe.ImagePath && <img src={recipe.ImagePath} alt={recipe.Name} />}
      <h2>Core Properties</h2>
      <ul>
        {Object.entries(recipe.CoreProps).map(([key, value]) => (
          <li key={key}>{key}: {value}</li>
        ))}
      </ul>
      <h2>Ingredients</h2>
      <ul>
        {recipe.Ingredients.map((ing, i) => <li key={i}>{ing}</li>)}
      </ul>
      <h2>Steps</h2>
      <ol>
        {recipe.Steps.map((step, i) => <li key={i}>{step}</li>)}
      </ol>
    </div>
  );
}
