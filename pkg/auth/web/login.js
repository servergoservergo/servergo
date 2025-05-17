document.addEventListener('DOMContentLoaded', function() {
    // 更新年份
    document.getElementById('current-year').textContent = new Date().getFullYear();
    
    // 获取登录表单和错误消息元素
    const loginForm = document.getElementById('login-form');
    const errorMessage = document.getElementById('error-message');
    
    // 从URL中获取错误消息（如果有）
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    
    if (error) {
        // 显示错误消息
        errorMessage.textContent = decodeURIComponent(error);
        errorMessage.style.display = 'block';
    }
    
    // 表单提交处理
    loginForm.addEventListener('submit', function(event) {
        // 获取用户名和密码
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value;
        
        // 获取国际化的错误消息
        const errorEmptyFields = window.loginErrorMessages.emptyFields;
        
        // 简单的客户端验证
        if (!username || !password) {
            event.preventDefault();
            errorMessage.textContent = errorEmptyFields;
            errorMessage.style.display = 'block';
            return;
        }
        
        // 如果验证通过，表单会正常提交到服务器
    });
    
    // 输入时隐藏错误消息
    document.getElementById('username').addEventListener('input', function() {
        errorMessage.style.display = 'none';
    });
    
    document.getElementById('password').addEventListener('input', function() {
        errorMessage.style.display = 'none';
    });
}); 