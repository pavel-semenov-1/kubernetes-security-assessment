import { useEffect, useState } from "react";
import { Search } from "lucide-react";

interface SearchInputProps {
    className: string;
    setLoading: (newValue: boolean) => void;
    setQuery: (query: string) => void;
};

// Custom search input textbox
const SearchInput = ({ className, setLoading, setQuery }: SearchInputProps) => {
  const [input, setInput] = useState("");
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  useEffect(() => {
    setLoading(true)
    const timer = setTimeout(() => {
      if (input) {
        setQuery(input);
      } else {
        setQuery("");
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [input]);

  return (
    <div className={className}>
      <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
      <input
        type="text"
        value={input}
        onChange={handleChange}
        placeholder="Search..."
        className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
      />
    </div>
  );
}

export default SearchInput;
