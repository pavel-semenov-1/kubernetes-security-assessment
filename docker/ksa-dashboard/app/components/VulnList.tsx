'use client';
import React, { useState } from "react";
import { ItemListProps } from "../types/ItemListProps";
import { Plus, Minus } from 'lucide-react';
import { Vulnerability } from "../types/Vulnerability";

const VulnList = ({ data }: ItemListProps<Vulnerability>) => {
    const [openCategory, setOpenCategory] = useState<number | null>(0);
    const [openItem, setOpenItem] = useState<number | null>(null);

    return (
        <div className="p-4 flex-auto">
            {Object.entries(data).map(([category, items], index) => (
                <div key={index} className="mb-4">
                    <button
                        className="flex flex-row w-full bg-blue-700 text-white font-semibold py-2 px-4 rounded-md text-left"
                        onClick={() =>
                            setOpenCategory(openCategory === index ? null : index)
                        }
                    >
                        <div className="grow flex flex-row">{openCategory === index ? (<Minus className="mr-2"/>) : (<Plus className="mr-2"/>)}{category}</div><p>{items.length }</p>
                    </button>

                    {openCategory === index && (
                        <ul className="mt-2 bg-gray-100 p-2 rounded-md shadow text-gray-900">
                            {items.length == 0 ?
                             (<div className="text-gray-500">No data</div>) : 
                             items.map((item, id) => (
                                <li
                                key={id}
                                className={`border-transparent border-2 rounded-md last:border-b-0 hover:border-blue-500 ${openItem === id && 'bg-slate-300'}`}
                            >
                                <a className="items-center block flex flex-row p-3" onClick={() =>
                                    setOpenItem(openItem === id ? null : id)
                                }>
                                    <span
                                        className={`transform transition-transform duration-300 ${
                                            openItem === id ? 'rotate-90' : ''
                                        }`}
                                        >
                                        <svg
                                            className="w-4 h-4"
                                            fill="none"
                                            stroke="currentColor"
                                            strokeWidth="2"
                                            viewBox="0 0 24 24"
                                            xmlns="http://www.w3.org/2000/svg"
                                        >
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7"></path>
                                        </svg>
                                    </span>
                                    <span className="font-semibold">{item?.ID}</span>: {item.Title}
                                </a>
                                {openItem === id && (
                                    <div className="text-gray-900 flex flex-col mt-1 px-7 pb-2">
                                        <div><p className="font-semibold inline">Target:</p> {item?.Target}</div>
                                        <div><p className="font-semibold inline">Package Name:</p> {item?.PkgName}</div>
                                        <div><p className="font-semibold inline">Description:</p> {item?.Description}</div>
                                        <div><p className="font-semibold inline">Installed Version:</p> {item?.InstalledVersion}</div>
                                        <div><p className="font-semibold inline">Fixed Version:</p> {item.FixedVersion || "unknown"}</div>
                                    </div>
                                )}
                            </li>
                            ))}
                        </ul>
                    )}
                </div>
            ))}
        </div>
    );
}

export default VulnList;