'use client'
import React, {useState} from 'react';
import SubjectTable from './components/SubjectTable';
import SubjectInput from './components/SubjectInput';

export default function Home() {
  const [refresh, setRefresh] = useState(false);

  const handleFormSubmit = () => {
    setRefresh(!refresh);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm lg:flex">
        <div>
          <h1>Subject List</h1>
          <SubjectTable refresh={refresh}/>
        </div>
        <div>
          <h1>New Subject</h1>
          <SubjectInput onFormSubmit={handleFormSubmit}/>
        </div>
      </div>
    </main>
  );
}
