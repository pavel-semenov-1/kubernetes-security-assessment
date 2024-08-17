'use client'
import React, { useEffect, useState } from 'react';
import { User } from '../types/User';
import { Subject } from '../types/Subject';  // Assuming you have a Subject type

const UserTable = ({ refresh }: { refresh: any }) => {
  const [users, setUsers] = useState<User[]>([]);
  const [subjectsByUser, setSubjectsByUser] = useState<Map<number, Subject[]>>(new Map());
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch(`/api/users`);
        if (!response.ok) {
          throw new Error('Failed to fetch users');
        }
        const data: User[] = await response.json();
        setUsers(data);

        const subjectsMap = new Map<number, Subject[]>();
        for (const user of data) {
          const subjectsResponse = await fetch(`/api/users/${user.id}/students`);
          if (!subjectsResponse.ok) {
            throw new Error('Failed to fetch subjects for user');
          }
          const subjectsData = await subjectsResponse.json();
          
          const subjectDetailsPromises = subjectsData.map((subj: { subjectId: number }) =>
            fetch(`/api/subjects/${subj.subjectId}`).then(res => res.json())
          );
          const subjectsDetails = await Promise.all(subjectDetailsPromises);

          subjectsMap.set(user.id, subjectsDetails);
        }
        setSubjectsByUser(subjectsMap);

      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
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
          <th className="text-sm font-medium px-6 py-4 text-left">Email</th>
          <th className="text-sm font-medium px-6 py-4 text-left">Subjects</th>
        </tr>
      </thead>
      <tbody>
        {users.map(user => (
          <tr key={user.id} className="border-b">
            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">{user.id}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">{user.name}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">{user.email}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">
              {subjectsByUser.get(user.id)?.map((subject: any) => (
                <div key={subject.id}>{subject.name}</div>
              )) || 'No subjects'}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default UserTable;
