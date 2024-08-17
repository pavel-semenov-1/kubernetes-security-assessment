'use client'
import React, {useState} from 'react';
import UserTable from './components/UserTable';
import UserInput from './components/UserInput';

export default function Home() {
  const [refresh, setRefresh] = useState(false);

  const handleFormSubmit = () => {
    setRefresh(!refresh);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm lg:flex">
        <div>
          <h1>User List</h1>
          <UserTable refresh={refresh}/>
        </div>
        <div>
          <h1>New User</h1>
          <UserInput onFormSubmit={handleFormSubmit}/>
        </div>
      </div>
    </main>
  );
}
