import { Component, lazy } from 'solid-js';
import { Router, Route, Navigate} from '@solidjs/router';
import { isAuthenticated } from './services/AuthServices';

import AuthGuard from './components/AuthGuard';
import AdminPanel from './components/AdminPanel';
import UserLayout from './layouts/UserLayout';

import './styles/global.css';

const TermoDeAceite = lazy(() => import('./components/TermoDeAceite'));
const HomePage = lazy(() => import('./layouts/UserLayout'));
const AdminLogin = lazy(() => import('./components/LoginAdmin'));
const AdminRegister = lazy(() => import('./components/RegisterAdmin'));
const AdminPage = lazy(() => import('./components/AdminPanel'));

const App: Component = () => {
  return (
    <Router>
        <Route path="/" component={TermoDeAceite} />

        <Route path="/sendMsg" component={UserLayout} />
        
        <Route path="/login" component={() => 
            isAuthenticated() ? <Navigate href="/admin" /> : <AdminLogin />
        } />
        <Route path="/register" component={() => 
          isAuthenticated() ? <Navigate href="/admin" /> : <AdminRegister />
        } />
        
        <Route path="/admin" component={() => 
          <AuthGuard>
            <AdminPage />
          </AuthGuard>
        } />
        <Route path="/admin/*" component={() => 
          <AuthGuard>
            <AdminPage />
          </AuthGuard>
        } />
    </Router>
  );
};

export default App;