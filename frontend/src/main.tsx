import ReactDOM from 'react-dom/client';
import {BrowserRouter as Router} from 'react-router-dom'; // 导入 Router
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
    <Router> {/* 使用 Router 包裹 App */}
        <App/>
    </Router>
);
