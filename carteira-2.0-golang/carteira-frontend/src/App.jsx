import { useState } from 'react';
import Portfolio_Screen from './components/Portfolio_Screen';
import Analysis_Screen from './components/Analysis_Screen';

function App() {
  const [currentScreen, setCurrentScreen] = useState('portfolio'); // 'portfolio' | 'analysis'

  return (
    <>
      {currentScreen === 'portfolio' && (
        <Portfolio_Screen onNavigate={setCurrentScreen} />
      )}
      {currentScreen === 'analysis' && (
        <Analysis_Screen onNavigate={setCurrentScreen} />
      )}
    </>
  );
}

export default App;
