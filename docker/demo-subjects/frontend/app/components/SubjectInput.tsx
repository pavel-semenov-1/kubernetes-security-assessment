'use client';
import React, { useState, useEffect } from 'react';
import { Subject } from '../types/Subject';
import { User } from '../types/User';

const SubjectInput = ({ onFormSubmit }: { onFormSubmit: () => void }) => {
  const [subject, setSubject] = useState<Subject>({ id: 0, name: '' });
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [students, setStudents] = useState<User[]>([]);
  
  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch('/api/users');
        const data = await response.json();
        setUsers(data);
      } catch (error) {
        console.error('Failed to fetch users:', error);
      }
    };
    fetchUsers();
  }, []);

  const handleAddUser = () => {
    if (selectedUserId !== null) {
      const userToAdd = users.find((user: User) => user.id === selectedUserId);
      if (userToAdd && !students.some((user: User) => user.id === userToAdd.id)) {
        setStudents((prev: User[]) => [...prev, userToAdd]);
      }
    }
  };

  const action = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
        if (subject.name.trim().length > 0) {
            const subjectResponse = await fetch('/api/subjects', {
                method: 'post',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: subject.name }),
            });
            const createdSubject = await subjectResponse.json();

            for (const user of students) {
                await fetch(`/api/students`, {
                    method: 'post',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        subjectId: createdSubject.id,
                        userId: user.id,
                    }),
                });
            }

            onFormSubmit();
        }
    } catch (error) {
        console.error('Error during submission:', error);
    }
  };

  return (
    <form className="form" onSubmit={action}>
      <label className="my-2 block text-sm font-medium text-gray-700 dark:text-gray-100">
        Name
        <div className="relative mt-1">
          <input
            type="text"
            className="block w-full h-10 pl-8 pr-3 mt-1 text-sm text-gray-700 border focus:outline-none rounded shadow-sm focus:border-blue-500"
            placeholder="Subject Name"
            value={subject.name}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSubject({ ...subject, name: e.target.value })}
          />
        </div>
      </label>
      <label className="my-2 block text-sm font-medium text-gray-700 dark:text-gray-100">
        Students
        <div className="relative mt-1">
          <select
            className="block w-full h-10 pl-3 pr-8 mt-1 text-sm text-gray-700 border focus:outline-none rounded shadow-sm focus:border-blue-500"
            value={selectedUserId || ''}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => setSelectedUserId(parseInt(e.target.value))}
          >
            <option value="" disabled>Select a user</option>
            {users.map((user: User) => (
              <option key={user.id} value={user.id}>
                {user.name}
              </option>
            ))}
          </select>
          <button
            type="button"
            className="my-2 px-4 py-1 min-w-[100px] w-full text-center text-white bg-green-600 border border-green-600 rounded active:text-green-500 hover:bg-transparent hover:text-green-600 focus:outline-none focus:ring"
            onClick={handleAddUser}
          >
            Add
          </button>
        </div>
        <div className="mt-2">
          {students.length > 0 && (
            <ul className="list-disc pl-5">
              {students.map((student: User) => (
                <li key={student.id}>{student.name}</li>
              ))}
            </ul>
          )}
        </div>
      </label>
      <input
        className="my-2 px-6 py-1 min-w-[100px] w-full text-center text-white bg-violet-600 border border-violet-600 rounded active:text-violet-500 hover:bg-transparent hover:text-violet-600 focus:outline-none focus:ring"
        type="submit"
        value="Submit"
      />
    </form>
  );
};

export default SubjectInput;
