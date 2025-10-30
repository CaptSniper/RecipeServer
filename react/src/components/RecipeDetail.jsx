import React, { useEffect, useState } from 'react';
import { fetchRecipe } from '../api/recipes';
import { useParams } from 'react-router-dom';
import './RecipeDetail.css';

export default function RecipeDetail() {
  const { id } = useParams();
  const [recipe, setRecipe] = useState(null);

  useEffect(() => {
    fetchRecipe(id).then(setRecipe);
  }, [id]);

  if (!recipe) return (
    <div className="recipe-detail-container">
      <div className="loading">Loading recipe...</div>
    </div>
  );

  return (
    <article className="recipe-detail-container">
      {/* Recipe Header */}
      <header className="recipe-header">
        <h1 className="recipe-title">{recipe.Name}</h1>

        {/* Core Properties (Cook Time, Prep Time, etc.) */}
        {recipe.CoreProps && Object.keys(recipe.CoreProps).length > 0 && (
          <div className="recipe-meta">
            {Object.entries(recipe.CoreProps).map(([key, value]) => (
              <div key={key} className="meta-item">
                <span className="meta-label">{key}</span>
                <span className="meta-value">{value}</span>
              </div>
            ))}
          </div>
        )}
      </header>

      {/* Recipe Image */}
      {recipe.ImagePath && (
        <figure className="recipe-image-container">
          <img
            src={recipe.ImagePath}
            alt={recipe.Name}
            className="recipe-image"
          />
        </figure>
      )}

      {/* Ingredients Section */}
      <section className="recipe-section">
        <h2 className="section-title">Ingredients</h2>
        <ul className="ingredients-list">
          {recipe.Ingredients.map((ing, i) => (
            <li key={i} className="ingredient-item">
              <span className="ingredient-bullet">•</span>
              <span className="ingredient-text">{ing}</span>
            </li>
          ))}
        </ul>
      </section>

      {/* Steps Section */}
      <section className="recipe-section">
        <h2 className="section-title">Preparation</h2>
        <ol className="steps-list">
          {recipe.Steps.map((step, i) => (
            <li key={i} className="step-item">
              <span className="step-number">{i + 1}</span>
              <p className="step-text">{step}</p>
            </li>
          ))}
        </ol>
      </section>
    </article>
  );
}
