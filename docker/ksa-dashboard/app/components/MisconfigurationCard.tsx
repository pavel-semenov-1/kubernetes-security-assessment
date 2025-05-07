'use client'
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { Misconfiguration } from "../types/Misconfiguration";
import { CircleX, CircleCheck, Trash2 } from "lucide-react";
import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { ITEM_STATUS_PASS, ITEM_STATUS_RESOLVED } from "../constants";

interface MisconfigurationCardProps {
    key: number;
    id: number;
    item: Misconfiguration;
    openItem: number | null;
    setOpenItem: (value: number | null) => void;
    onClickResolve: () => void;
    onClickDelete: () => void;
}

const MisconfigurationCard = ({ id, item, openItem, setOpenItem, onClickDelete, onClickResolve }: MisconfigurationCardProps) => {
    const maximumTitleLength = 165;
    const [opaque, setOpaque] = useState<boolean>(item.Status == ITEM_STATUS_RESOLVED || item.Status == ITEM_STATUS_PASS);
    useEffect(() => {
        setOpaque(item.Status === ITEM_STATUS_RESOLVED || item.Status === ITEM_STATUS_PASS);
    }, [item.Status]);
    const truncated = (text: string, maxLength: number) => 
        text.length > maxLength ? text.slice(0, maxLength) + "..." : text;

    return (
        <motion.li
            key={id}
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: opaque ? 0.5 : 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.9, height: 0, margin: 0, padding: 0 }}
            transition={{ duration: 0.3 }}
            className={`flex flex-col w-full px-4 border-2 border-transparent rounded-md 
                hover:border-blue-500 
                last:border-b-2 last:border-transparent last:hover:border-blue-500 
                ${openItem === id && 'bg-slate-300'}`}
        >
            <a
                className="flex flex-row items-center block py-3 cursor-pointer"
                onClick={() => setOpenItem(openItem === id ? null : id)}
            >
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
                <span
                    className={`ml-2 font-semibold transition-all ${
                        opaque ? 'line-through text-muted-foreground' : ''
                    }`}
                >
                    {truncated(item?.Target, maximumTitleLength - item.Title.length)}
                </span>
                : {item.Title}
            </a>

            {openItem === id && (
                <div className="text-gray-900 flex flex-col mt-1 px-4 pb-2">
                    <div>
                        <p className="font-semibold inline">Misconfiguration ID:</p> {item?.MisconfigurationID}
                    </div>
                    <div>
                        <p className="font-semibold inline">Type:</p> {item?.Type}
                    </div>
                    <div>
                        <p className="font-semibold inline">Description:</p> {item?.Description}
                    </div>
                    <div>
                        <p className="font-semibold inline">Resolution:</p> {item?.Resolution || "No resolution available."}
                    </div>
                    <div>
                        <p className="font-semibold inline">Status:</p> {item?.Status}
                    </div>
                    <div className="my-3 flex flex-row-reverse gap-2">
                        { item.Status !== ITEM_STATUS_PASS ? <Tooltip>
                            <TooltipTrigger asChild>
                                <Button
                                    onClick={() => {onClickResolve(); setOpaque(!opaque)}}
                                    className="w-10 bg-green-800 hover:bg-green-600 text-white"
                                    variant="ghost"
                                    size="icon"
                                >
                                    {opaque ? <CircleX size={18} /> : <CircleCheck size={18} />}
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                {opaque ? "Reopen" : "Resolve"}
                            </TooltipContent>
                        </Tooltip> : <></>}

                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Button
                                    onClick={onClickDelete}
                                    className="w-10 bg-red-800 hover:bg-red-600 text-white"
                                    variant="ghost"
                                    size="icon"
                                >
                                    <Trash2 size={18} />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>Delete</TooltipContent>
                        </Tooltip>
                    </div>
                </div>
            )}
        </motion.li>
    );
};

export default MisconfigurationCard;