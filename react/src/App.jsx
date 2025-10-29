import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RecipeList from './components/RecipeList';
import RecipeDetail from './components/RecipeDetail';
import RecipeForm from './components/RecipeForm';
import ScrapeRecipe from './pages/ScrapeRecipe'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<RecipeList />} />
        <Route path="/recipes" element={<RecipeList />} />
        <Route path="/recipe/:id" element={<RecipeDetail />} />
        <Route path="/new" element={<RecipeForm />} />
        <Route path="/scrape" element={<ScrapeRecipe />} />
        <Route path="/scrape/:id" element={<ScrapeRecipe />} /> {/* for editing existing recipe */}
      </Routes>
    </Router>
  );
}

export default App;
