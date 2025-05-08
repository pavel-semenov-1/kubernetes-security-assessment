'use client';
import React, { useEffect, useState } from "react";
import { Report } from "../types/Report";
import { RotateCcw, CloudDownload, Trash2, CheckCheck } from "lucide-react";
import SearchInput from "./SearchInput";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Checkbox } from "@/components/ui/checkbox";
import { ALL_SCANNERS, ITEM_TYPE_MISCONFIGURATION, ITEM_TYPE_VULNERABILITY } from "../constants";

interface FilterPanelProps {
    valueSelectedScanner: string;
    onChangeSelectedScanner: (newValue: string) => void;
    valueSelectedType: string;
    onChangeSelectedType: (newValue: string) => void;
    valueSelectedReport: string;
    onChangeSelectedReport: (newValue: string) => void;
    valueReports: Report[];
    showResolved: boolean;
    setShowResolved: (newValue: boolean) => void;
    onClickRescan: () => void;
    jobInProgress: boolean;
    setLoading: (newValue: boolean) => void;
    onChangeSearchValue: (newValue: string) => void;
    onClickDownload: () => void;
    onClickDelete: () => void;
    onClickResolve: () => void;
};

// A panel with filters and search on the top of the page
const FilterPanel = ({ valueSelectedScanner, onChangeSelectedScanner, valueSelectedType, onChangeSelectedType, valueSelectedReport, onChangeSelectedReport, valueReports, showResolved, setShowResolved, onClickRescan, jobInProgress, setLoading, onChangeSearchValue, onClickDownload, onClickDelete, onClickResolve }: FilterPanelProps) => {
    const types = [ITEM_TYPE_MISCONFIGURATION, ITEM_TYPE_VULNERABILITY];
    const [scanners, setScanners] = useState<string[]>([]);
    const [refreshDisabled, setRefreshDisabled] = useState<boolean>(false);

    const vulnerabilitySelected = () => {
        return valueSelectedType == ITEM_TYPE_VULNERABILITY
    }
    const disableButtons = () => {
        return valueSelectedReport === null || valueReports === null || valueReports.length == 0
    }
    const disableMisconfigButtons = () => {
        return disableButtons() || vulnerabilitySelected();
    }

    useEffect(() => {
        const fetchScanners = async () => {
            try {
                const response = await fetch(`/api/scanners`)
                if (!response.ok) throw new Error(`Failed to fetch scanners`)
                const data: string[] = await response.json()
                if (data === null || Object.keys(data).length === 0) {
                    setScanners([ALL_SCANNERS])
                } else {
                    setScanners([ALL_SCANNERS, ...data])
                }
            } catch (err: any) {
                console.error(err)
            } 
        }
        fetchScanners()
    }, [])

    return (
        <div className="flex items-center gap-4 p-4">
            <div className="flex items-center gap-3">
                <label htmlFor="scanner" className="text-base text-gray-800">
                    Scanner:
                </label>
                <select
                    id="scanner"
                    value={valueSelectedScanner}
                    onChange={(e) => {onChangeSelectedScanner(e.target.value); setRefreshDisabled(e.target.value === scanners[0])}}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                >
                    {scanners.map((scanner, index) => (
                        <option key={index} value={scanner}>
                            {scanner}
                        </option>
                    ))}
                </select>
            </div>

            <div className="flex items-center gap-3">
                <label htmlFor="type" className="text-base text-gray-800">
                    Type:
                </label>
                <select
                    id="type"
                    value={valueSelectedType}
                    onChange={(e) => onChangeSelectedType(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {types.map((type, index) => (
                        <option key={index} value={type}>
                        {type}
                        </option>
                    ))}
                    </select>
            </div>

            <div className="flex items-center gap-3">
                <label htmlFor="report" className="text-base text-gray-800">
                    Report:
                </label>
                {valueReports !== null && valueReports.length > 0 ?
                <select
                    id="report"
                    value={valueSelectedReport}
                    onChange={(e) => onChangeSelectedReport(e.target.value)}
                    className="max-w-[15rem] truncate text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {valueReports.map((report) => (
                        <option key={report.ID} value={report.ID}>
                        {report.Filename}
                        </option>
                    ))}
                </select>
                : valueSelectedScanner === "All scanners" ? <span>Newest reports</span> :
                <span>No data. Trigger a scan first.</span>
                }
            </div>
            <div className={`group flex flex-row gap-1 border-2 rounded-2xl ${vulnerabilitySelected() ? "border-gray-500" : "border-blue-700"} bg-white px-3 py-[6px] ${vulnerabilitySelected() ? "" : "hover:bg-blue-700"}`}>
                <Checkbox
                    disabled={vulnerabilitySelected()}
                    checked={showResolved}
                    onCheckedChange={setShowResolved}
                    id="showResolved"
                    className={`border ${vulnerabilitySelected() ? "border-gray-500" : "border-blue-700 group-hover:border-white"} group-hover:bg-white 
                                text-blue-700 data-[state=checked]:bg-blue-700 
                                data-[state=checked]:text-white 
                                group-hover:data-[state=checked]:text-blue-700
                                transition-colors`}
                    />
                <label
                    htmlFor="showResolved"
                    className={`text-sm font-medium leading-none text-black ${vulnerabilitySelected() ? "" : "group-hover:text-white"} peer-disabled:cursor-not-allowed peer-disabled:opacity-70 transition-colors`}
                >
                    Show resolved
                </label>
            </div>

            <div className="flex flex-row gap-2">
                <div className={( refreshDisabled ? "bg-gray-500" : "bg-blue-700") + " text-white border rounded-full p-1"}  style={{ pointerEvents: refreshDisabled ? "none" : "auto"}}>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <a onClick={ refreshDisabled ? undefined : onClickRescan} className="button"><RotateCcw className={jobInProgress ? "animate-spin" : ""}/></a>
                        </TooltipTrigger>
                        <TooltipContent>
                            Rescan
                        </TooltipContent>
                    </Tooltip>
                </div>
                <div className={`border rounded-full p-1 ${disableButtons() ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-700'}`}>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <a
                                className={`button text-white ${disableButtons() ? 'pointer-events-none opacity-50' : ''}`}
                                onClick={disableButtons() ? undefined : onClickDownload}
                            >
                                <CloudDownload/>
                            </a>
                        </TooltipTrigger>
                        <TooltipContent>
                            Download report
                        </TooltipContent>
                    </Tooltip>
                </div>
                <div className={`border rounded-full p-1 ${disableButtons() ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-700'}`}>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <a
                                className={`button text-white ${disableButtons() ? 'pointer-events-none opacity-50' : ''}`}
                                onClick={disableButtons() ? undefined : onClickDelete}
                            >
                                <Trash2/>
                            </a>
                        </TooltipTrigger>
                        <TooltipContent>
                            Delete filtered
                        </TooltipContent>
                    </Tooltip>
                </div>
                <div className={`border rounded-full p-1 ${disableMisconfigButtons() ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-700'}`}>
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <a
                                className={`button text-white ${disableButtons() ? 'pointer-events-none opacity-50' : ''}`}
                                onClick={disableMisconfigButtons() ? undefined : onClickResolve}
                            >
                                <CheckCheck />
                            </a>
                        </TooltipTrigger>
                        <TooltipContent>
                            Resolve filtered
                        </TooltipContent>
                    </Tooltip>
                </div>
            </div>

            <SearchInput setLoading={setLoading} setQuery={onChangeSearchValue} className="relative ml-auto"/>
        </div>
    );
}

export default FilterPanel;