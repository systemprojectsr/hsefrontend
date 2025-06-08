const API_BASE = 'http://176.57.215.221:8080/';
const form = document.querySelector('.companyRegistrationForm');

async function handleCompanyRegister(event) {
  event.preventDefault();
  const body = {
    user: {
      register: {
        company_name: document.querySelector('.companyName').value,
        email: document.querySelector('.companyEmail').value,
        phone: document.querySelector('.companyPhone').value,
        full_name: document.querySelector('.agentName').value,
        position_agent: document.querySelector('.agentPosition').value,
        id_company: document.querySelector('.companyId').value,
        address: document.querySelector('.companyAddress').value,
        type_service: document.querySelector('.serviceType').value,
        password: document.querySelector('.companyPassword').value,
        photo: null,
        documents: [],
        type: 'company'
      }
    }
  };
  try {
    const response = await fetch(`${API_BASE}v1/register/company`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    });
    const data = await response.json();
    if (response.ok) {
      const { token, type } = data.user;
      localStorage.setItem('authToken', token);
      localStorage.setItem('accountType', type);
      document.cookie = `authToken=${token}; path=/`;
      window.location.href = 'companyProfile.html';
    } else {
      console.error('Registration error');
    }
  } catch (err) {
    console.error('Ошибка:', err);
  }
}

if (form) {
  form.addEventListener('submit', handleCompanyRegister);
}
