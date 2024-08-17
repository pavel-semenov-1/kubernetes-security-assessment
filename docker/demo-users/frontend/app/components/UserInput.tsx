'use client'
import React, { useState } from 'react';
import { User } from '../types/User';

const UserInput  = ({ onFormSubmit }: { onFormSubmit: () => void}) => {
  const [user, setUser] = useState<User>({id: 0, name: '', email: ''})
  const action = async (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
        if (user.name.trim().length > 0 && user.email.trim().length > 0) {
            await fetch(`/api/users`, {
                method: 'post',
                headers: {'Content-Type':'application/json'},
                body: JSON.stringify(user)
            });
            onFormSubmit();
        }
    } catch (error) {
        console.error(error);
    }
  }

  return (
    <form className="form" onSubmit={action}>
        <label className="my-2 block text-sm font-medium text-gray-700 dark:text-gray-100">
        Name
            <div className="relative mt-1">
                <input type="text" id="input-6" className="block w-full h-10 pl-8 pr-3 mt-1 text-sm text-gray-700 border focus:outline-none rounded shadow-sm focus:border-blue-500" placeholder="John Doe" name="name" value={user.name} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setUser({id: user.id, name: e.target.value, email: user.email})}/>
                <span className="absolute inset-y-0 left-0 flex items-center justify-center ml-2"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className="w-4 h-4 text-blue-400 pointer-events-none">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zM6 20a6 6 0 0112 0H6z"></path>
                </svg>
                </span>
            </div>      
        </label>
        <label className="my-2 block text-sm font-medium text-gray-700 dark:text-gray-100">
        Email
            <div className="relative mt-1">
                <input type="email" id="input-6" className="block w-full h-10 pl-8 pr-3 mt-1 text-sm text-gray-700 border focus:outline-none rounded shadow-sm focus:border-blue-500" placeholder="user@xyz.com" name="email" value={user.email} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setUser({id: user.id, name: user.name, email: e.target.value})}/>
                <span className="absolute inset-y-0 left-0 flex items-center justify-center ml-2"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className="w-4 h-4 text-blue-400 pointer-events-none"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
                </span>
            </div>        
        </label>
        <input className="my-2 px-6 py-1 min-w-[100px] w-full text-center text-white bg-violet-600 border border-violet-600 rounded active:text-violet-500 hover:bg-transparent hover:text-violet-600 focus:outline-none focus:ring" type="submit" value="Submit" />
    </form>
  );
};

export default UserInput;