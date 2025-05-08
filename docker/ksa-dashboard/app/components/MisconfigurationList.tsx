'use client';
import React from "react";
import { ItemListProps } from "../types/ItemListProps";
import { Plus, Minus } from 'lucide-react';
import { Misconfiguration } from "../types/Misconfiguration";
import MisconfigurationCard from "./MisconfigurationCard";
import { AnimatePresence } from "framer-motion";
import { ITEM_STATUS_FAIL, ITEM_STATUS_MANUAL } from "../constants";

// List of misconfigurations
const MisconfigurationList = ({ data, onClickDelete, onClickResolve, openCategory, setOpenCategory, openItem, setOpenItem }: ItemListProps<Misconfiguration>) => {
    const emptyDataMessage = "No data. Run a new scan or a select a different object type.";

    return (
        <div className="p-4 flex-auto">
            {data === null ? emptyDataMessage : Object.entries(data).map(([category, items], index) => (
                <div key={index} className="mb-4">
                    <button
                        className="flex flex-row w-full bg-blue-700 text-white font-semibold py-2 px-4 rounded-md text-left"
                        onClick={() =>
                            setOpenCategory(openCategory === index ? null : index)
                        }
                    >
                        <div className="grow flex flex-row">{openCategory === index ? (<Minus className="mr-2"/>) : (<Plus className="mr-2"/>)}{category}</div>
                        <p>{items.filter(i => i.Status == ITEM_STATUS_FAIL || i.Status == ITEM_STATUS_MANUAL).length}/{items.length}</p>
                    </button>

                    {openCategory === index && (
                        <ul className="mt-2 bg-gray-100 p-2 rounded-md shadow text-gray-900">
                            <AnimatePresence>
                                {items.length == 0 ?
                                (<div className="text-gray-500">No data</div>) : 
                                items.map((item, id) => <MisconfigurationCard key={id} id={id} item={item} openItem={openItem} setOpenItem={setOpenItem} onClickDelete={() => 
                                 onClickDelete(category, id, item)} onClickResolve={() => onClickResolve ? onClickResolve(category, id, item) : {}}/>)}
                            </AnimatePresence>
                        </ul>
                    )}
                </div>
            ))}
        </div>
    );
}

export default MisconfigurationList;