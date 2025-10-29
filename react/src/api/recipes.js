import axios from 'axios';

const BASE_URL = '/api'; // Go server API

export async function fetchRecipes() {
  try {
    const res = await axios.get(`${BASE_URL}/recipes`);
    const data = res.data;
    // Always return an array
    if (Array.isArray(data)) return data;
    if (Array.isArray(data.recipes)) return data.recipes;
    return [];
  } catch (err) {
    console.error('Failed to fetch recipes:', err);
    return [];
  }
}

export async function fetchRecipe(id) {
  try {
    const res = await axios.get(`${BASE_URL}/recipes/${id}`);
    return res.data;
  } catch (err) {
    console.error(`Failed to fetch recipe ${id}:`, err);
    throw err;
  }
}

export async function createRecipe(recipe) {
  try {
    const res = await axios.post(`${BASE_URL}/recipes`, recipe);
    return res.data;
  } catch (err) {
    console.error('Failed to create recipe:', err);
    throw err;
  }
}

export async function updateRecipe(id, recipe) {
  try {
    const res = await axios.put(`${BASE_URL}/recipes/${id}`, recipe);
    return res.data;
  } catch (err) {
    console.error(`Failed to update recipe ${id}:`, err);
    throw err;
  }
}

export async function deleteRecipe(id) {
  try {
    const res = await axios.delete(`${BASE_URL}/recipes/${id}`);
    return res.data;
  } catch (err) {
    console.error(`Failed to delete recipe ${id}:`, err);
    throw err;
  }
}

export async function scrapeRecipe(url) {
  try {
    const res = await axios.post('/api/scrape', { url });
    return res.data;
  } catch (err) {
    console.error('Failed to scrape recipe:', err);
    throw err;
  }
}
