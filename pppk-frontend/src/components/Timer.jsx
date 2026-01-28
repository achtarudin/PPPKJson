import React, { useEffect, useState } from 'react';

const Timer = ({ expiresAt, onExpire }) => {
  const [remaining, setRemaining] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      const diff = Math.max(0, Math.floor((new Date(expiresAt) - new Date()) / 1000));
      setRemaining(diff);
      if (diff === 0) {
        clearInterval(interval);
        if (onExpire) onExpire();
      }
    }, 1000);
    return () => clearInterval(interval);
  }, [expiresAt, onExpire]);

  const minutes = Math.floor(remaining / 60);
  const seconds = remaining % 60;
  let color = 'text-success';
  if (remaining <= 600) color = 'text-danger';
  else if (remaining <= 1800) color = 'text-warning';

  return (
    <span className={`fw-bold ${color}`}>{minutes}:{seconds.toString().padStart(2, '0')}</span>
  );
};

export default Timer;
