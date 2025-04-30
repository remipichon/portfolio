import React, {useEffect, useState} from "react";


type Job = {
    namespace: string;
    name: string;
    lastSuccessfullyRunStarTime?: Date;
    lastStatus: {
        type: string;
        message?: string;
    };
    lastSuccessfullyRunCompletionTime?: Date;
};

export function App() {
    const [jobs, setJobs] = useState<Job[]>([]);
    const [loading, setLoading] = useState(true);
    const [lastFetchJobs, setLastFetchJobs] = useState<Date | null>(null);

    const fetchJobs = async () => {
        setLoading(true);

        const jobsRaw = await fetch("/list").then(res => res.json());
        const jobs: Job[] = jobsRaw["jobs"].map(parseJob);
        setJobs(jobs);
        setLoading(false);
        setLastFetchJobs(new Date());
    };

    const parseJob = (raw: any): Job => ({
        ...raw,
        lastSuccessfullyRunStarTime: raw.lastSuccessfullyRunStarTime ? new Date(raw.lastSuccessfullyRunStarTime) : undefined,
        lastSuccessfullyRunCompletionTime: raw.lastSuccessfullyRunCompletionTime ? new Date(raw.lastSuccessfullyRunCompletionTime) : undefined,
    });

    useEffect(() => {
        // Fetch initially
        fetchJobs();

        // Set up visibility change listener
        const handleVisibilityChange = () => {
            if (document.visibilityState === 'visible') {
                console.info("Tab just got visible, perform reload");
                fetchJobs();
            }
        };

        console.info("Add visibilitychange event listener");
        document.addEventListener('visibilitychange', handleVisibilityChange);

        // Set up polling
        const intervalId = setInterval(() => {
            if (document.visibilityState === 'visible') {
                console.info("Polling fetchJobs every 5s");
                fetchJobs();
            }
        }, 5000);

        // Cleanup on unmount
        return () => {
            console.info("Disable Polling fetchJobs every 5s");
            document.removeEventListener('visibilitychange', handleVisibilityChange);
            clearInterval(intervalId);
        };
    }, []);

    // refresh when the tab get the focus
    const handleVisibilityChange = () => {
        if (document.visibilityState === 'visible') {
            console.log("Tab just got visible, perform reload")
            fetchJobs()
        }
    };



    const runJob = async (namespace: string, name: string) => {
        await fetch(`/run/${namespace}/${name}`);
        fetchJobs();
    };

    const killJob = async (namespace: string, name: string) => {
        await fetch(`/kill/${namespace}/${name}`);
        fetchJobs();
    };

    return (
        <div style={{padding: "2rem", fontFamily: "Arial, sans-serif"}}>
            <h2>Job Assistant</h2>
            {lastFetchJobs && (
                <p>Last fetch: {lastFetchJobs.toLocaleTimeString()}</p>
            )}            {loading ? (
                <p>Loading jobs...</p>
            ) : (
                <table style={{width: "100%", borderCollapse: "collapse"}}>
                    <thead>
                    <tr>
                        <th style={thStyle}>Namespace</th>
                        <th style={thStyle}>Name</th>
                        <th style={thStyle}>Status</th>
                        <th style={thStyle}>Start time</th>
                        <th style={thStyle}>Completion time</th>
                        <th style={thStyle}>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {jobs.map((job) => (
                        <tr key={`${job.namespace}-${job.name}`}>
                            <td style={tdStyle}>{job.namespace}</td>
                            <td style={tdStyle}>{job.name}</td>
                            <td style={tdStyle}>{job.lastStatus?.type} - {job.lastStatus?.message}</td>
                            <td style={tdStyle}>{job.lastSuccessfullyRunStarTime?.toLocaleTimeString()}</td>
                            <td style={tdStyle}>{job.lastSuccessfullyRunCompletionTime?.toLocaleTimeString()}</td>
                            <td style={tdStyle}>
                                <button
                                    onClick={() => runJob(job.namespace, job.name)}
                                    disabled={job.lastStatus.type === "Running"}
                                    style={buttonStyle}
                                >
                                    Run
                                </button>
                                <button
                                    onClick={() => killJob(job.namespace, job.name)}
                                    disabled={job.lastStatus.type !== "Running"}
                                    style={{
                                        ...buttonStyle,
                                        marginLeft: "0.5rem"
                                    }}
                                >
                                    Kill
                                </button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}

const thStyle: React.CSSProperties = {
    textAlign: "left",
    padding: "8px",
    backgroundColor: "#f2f2f2",
    borderBottom: "1px solid #ccc",
};

const tdStyle: React.CSSProperties = {
    padding: "8px",
    borderBottom: "1px solid #eee",
};

const buttonStyle: React.CSSProperties = {
    padding: "6px 12px",
    fontSize: "14px",
    borderRadius: "4px",
    border: "1px solid #ccc",
    backgroundColor: "#f5f5f5",
    cursor: "pointer",
};

export default App;
