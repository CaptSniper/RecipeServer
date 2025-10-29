import React, { useState } from 'react';
import { createRecipe } from '../api/recipes';
import { useNavigate } from 'react-router-dom';

export default function RecipeForm() {
  const [name, setName] = useState('');
  const [ingredients, setIngredients] = useState(['']);
  const [steps, setSteps] = useState(['']);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    await createRecipe({ Name: name, Ingredients: ingredients.filter(Boolean), Steps: steps.filter(Boolean), CoreProps: {} });
    navigate('/');
  };

  return (
    <form onSubmit={handleSubmit}>
      <h1>New Recipe</h1>
      <input placeholder="Recipe Name" value={name} onChange={e => setName(e.target.value)} required />
      
      <h2>Ingredients</h2>
      {ingredients.map((ing, i) => (
        <input
          key={i}
          placeholder={`Ingredient ${i + 1}`}
          value={ing}
          onChange={e => {
            const copy = [...ingredients];
            copy[i] = e.target.value;
            setIngredients(copy);
          }}
        />
      ))}
      <button type="button" onClick={() => setIngredients([...ingredients, ''])}>Add Ingredient</button>

      <h2>Steps</h2>
      {steps.map((step, i) => (
        <input
          key={i}
          placeholder={`Step ${i + 1}`}
          value={step}
          onChange={e => {
            const copy = [...steps];
            copy[i] = e.target.value;
            setSteps(copy);
          }}
        />
      ))}
      <button type="button" onClick={() => setSteps([...steps, ''])}>Add Step</button>

      <br />
      <button type="submit">Create Recipe</button>
    </form>
  );
}
