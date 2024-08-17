'use client';
import React, { useEffect, useState } from 'react';
import { Subject } from '../types/Subject';
import { User } from '../types/User';

const SubjectTable = ({ refresh }: { refresh: any }) => {
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [studentsMap, setStudentsMap] = useState<Map<number, User[]>>(new Map());
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

  useEffect(() => {
    const fetchStudentsForSubjects = async () => {
      try {
        const studentsMap = new Map<number, User[]>();
        await Promise.all(
          subjects.map(async (subject: Subject) => {
            const response = await fetch(`/api/subjects/${subject.id}/students`);
            if (response.ok) {
              const studentData: { subjectId: number; userId: number }[] = await response.json();
              
              const userResponses = await Promise.all(
                studentData.map((student: any) => fetch(`/api/users/${student.userId}`))
              );
              const users = await Promise.all(userResponses.map((res) => res.json()));

              studentsMap.set(subject.id, users);
            }
          })
        );
        setStudentsMap(studentsMap);
      } catch (err: any) {
        setError('Failed to fetch students for subjects');
      }
    };

    if (subjects.length > 0) {
      fetchStudentsForSubjects();
    }
  }, [subjects]);

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
          <th className="text-sm font-medium px-6 py-4 text-left">Students</th>
        </tr>
      </thead>
      <tbody>
        {subjects.map((subject: Subject) => (
          <tr key={subject.id} className="border-b">
            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">{subject.id}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">{subject.name}</td>
            <td className="text-sm font-light px-6 py-4 whitespace-nowrap">
              {studentsMap.get(subject.id)
                ? studentsMap.get(subject.id)!.map((user: User) => user.name).join(', ')
                : 'No students'}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default SubjectTable;
