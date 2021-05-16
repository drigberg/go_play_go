import React from 'react';
import { BrowserRouter as Switch, Route } from 'react-router-dom';
import Home from './Home';

export default function AppRouter(): JSX.Element {
  return (
    <div>
      <Switch>
        <Route exact path="/" render={() => <Home />}></Route>
      </Switch>
    </div>
  );
}
