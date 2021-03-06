import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import AppRouter from './components/Router';

import './App.css';

export default function App(): JSX.Element {
  return (
    <Router>
      <AppRouter />
    </Router>
  );
}
