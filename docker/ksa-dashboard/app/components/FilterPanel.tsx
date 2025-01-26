'use client';
import React from "react";
import { FilterPanelProps } from "../types/FilterPanelProps";

const FilterPanel = ({ valueScanner, onChangeScanner, valueType, onChangeType }: FilterPanelProps) => {
    const scanners = ["Trivy", "Kubescape"];
    const types = ["Misconfiguration", "Vulnerability"];

    return (
        <div className="flex items-center gap-4 p-4 max-w-2xl">
            <div className="flex items-center gap-3">
                <label htmlFor="scanner" className="text-base text-gray-800">
                    Scanner:
                </label>
                <select
                    id="scanner"
                    value={valueScanner}
                    onChange={(e) => onChangeScanner(e.target.value)}
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
                    value={valueType}
                    onChange={(e) => onChangeType(e.target.value)}
                    className="text-gray-800 p-1 border rounded bg-white focus:outline-none focus:ring-2 focus:ring-darkblue-500"
                    >
                    {types.map((type, index) => (
                        <option key={index} value={type}>
                        {type}
                        </option>
                    ))}
                    </select>
            </div>
        </div>
    );
}

export default FilterPanel;