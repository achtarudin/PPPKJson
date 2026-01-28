import React from 'react';

const QuestionCard = ({ question, selectedOption, onSelect, questionNumber }) => {
  if (!question) return null;
  return (
    <div className="card">
      <div className="card-body">
        <div className="d-flex justify-content-between align-items-center mb-2">
          <span className="badge bg-primary">#{questionNumber}</span>
          <span className={`badge bg-$
            {question.category === 'MANAJERIAL' ? 'primary' :
              question.category === 'SOSIAL_KULTURAL' ? 'success' :
              question.category === 'TEKNIS' ? 'warning' :
              'info'}
          `}>{question.category}</span>
        </div>
        <h5 className="mb-3">{question.question_text}</h5>
        <form>
          {question.options.map(opt => (
            <div className="form-check mb-2" key={opt.id}>
              <input
                className="form-check-input"
                type="radio"
                name={`question-${question.exam_question_id}`}
                id={`option-${opt.id}`}
                checked={selectedOption === opt.id}
                onChange={() => onSelect(question.exam_question_id, opt.id)}
                disabled={typeof onSelect !== 'function'}
              />
              <label className="form-check-label" htmlFor={`option-${opt.id}`}>
                {opt.option_text}
              </label>
            </div>
          ))}
        </form>
      </div>
    </div>
  );
};

export default QuestionCard;
