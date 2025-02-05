import { FC, useState, useEffect } from 'react';
import { ChevronsRight, ChevronsLeft } from 'lucide-react';
import { StatResult } from '../types/StatResult';
import { Stats } from '../types/Stats';
import VulnerabilityBarChart from './VulnerabilityBarChart';
import MisconfigurationBarChart from './MisconfigurationBarChart';

const SidePanel: FC = () => {
    const [isPanelOpen, setIsPanelOpen] = useState<boolean>(false);
    const togglePanel = () => {
        setIsPanelOpen((prev) => !prev);
    };
    const [results, setResults] = useState<StatResult[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    // TODO: fetch scanners from the parser
    const scanners = ["trivy"]

    useEffect(() => {
      const fetchStats = async () => {
        try {
          const statsArray: StatResult[] = await Promise.all(
            scanners.map(async (scanner) => {
              const res = await fetch(`/api/stats?scanner=${scanner}`);
  
              if (!res.ok) {
                throw new Error(`Failed to fetch stats for ${scanner}: ${res.statusText}`);
              }
  
              const stats: Stats = await res.json();
              return { scanner, stats };
            })
          );
  
          setResults(statsArray);
        } catch (err) {
          setError(err instanceof Error ? err.message : "Unknown error");
        } finally {
          setLoading(false);
        }
      };
  
      fetchStats();
    }, []); 

    return (
    <div>
        <div
        className={`h-full w-2/5 fixed top-0 right-0 h-full bg-gray-600 text-white p-5 transform transition-all duration-300 ease-in-out ${
            isPanelOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
        >
            <div className={`${
              isPanelOpen ? 'left-[-35px]' : 'left-[-125px]'
            } gap-3 bg-blue-700 text-white rounded-md inline font-bold absolute top-4 bg-gray-600 hover:text-gray-300`}>
              <a
                className="p-2 block flex items-center button cursor-pointer transition-all"
                onClick={togglePanel}
              >
                {isPanelOpen ? <ChevronsLeft /> : <>STATISTICS<ChevronsRight /></>}
              </a>
            </div>
            {loading ? <div>Loading...</div> : error ? <div>{error}</div> :
                <>
                <div className="mb-2">
                    <h2 className="text-xl text-center">Vulnerabilities</h2>
                    <VulnerabilityBarChart data={results} />
                </div>
                <div>
                    <h2 className="text-xl text-center">Misconfigurations</h2>
                    <MisconfigurationBarChart data={results} />
                </div>
                </>
            }
        </div>
    </div>
    );
};

export default SidePanel;
