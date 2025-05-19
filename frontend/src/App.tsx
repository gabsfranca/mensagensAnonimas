import { Component } from 'solid-js';
import { Router, Route} from '@solidjs/router';

import { AdminPanel } from './components/AdminPanel';
import UserLayout from './layouts/UserLayout';

import './styles/global.css';

const App: Component = () => {
  return (
    <Router>
        <Route path="/" component={UserLayout} />
        <Route path="/admin" component={AdminPanel} />
    </Router>
  );
};

export default App;