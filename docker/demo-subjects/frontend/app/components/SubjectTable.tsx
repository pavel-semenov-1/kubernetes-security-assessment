'use client'
import React, { useEffect, useState } from 'react';
import { Subject } from '../types/Subject';

const SubjectTable = ({ refresh }: { refresh: any }) => {
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSubjects = async () => {
      try {
        const response = await fetch(`/api/subjects`);
        if (!response.ok) {
          throw new Error('Failed to fetch subjects');
        }
        const data: Subject[] = await response.json();
        setSubjects(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchSubjects();
  }, [refresh]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <table className="table min-w-full">
      <thead className="border-b">
        <tr>
          <th className="text-sm font-medium px-6 py-4 text-left">ID</th>
          <th className="text-sm font-medium px-6 py-4 text-left">Name</th>
        </tr>
      </thead>
      <tbody>
        {subjects.map(subject => (
          <tr key={subject.id} className="border-b">
            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">{subject.id}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">{subject.name}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default SubjectTable;